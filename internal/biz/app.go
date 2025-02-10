package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/conf"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli"
	helmValues "helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	helmrepo "helm.sh/helm/v3/pkg/repo"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

const (
	AppPackage string = "apps"

	SystemNamespace string = "cloud-copilot"

	DaemonSet   string = "DaemonSet"
	Deployment  string = "Deployment"
	StatefulSet string = "StatefulSet"
	ReplicaSet  string = "ReplicaSet"
	CronJob     string = "CronJob"
	Job         string = "Job"
	Pod         string = "Pod"
)

type AppUsecase struct {
	log  *log.Helper
	conf *conf.Bootstrap
}

func NewAppUseCase(conf *conf.Bootstrap, logger log.Logger) *AppUsecase {
	a := &AppUsecase{
		conf: conf,
		log:  log.NewHelper(logger),
	}
	return a
}

func (a *App) AddVersion(version *AppVersion) {
	if a.Versions == nil {
		a.Versions = make([]*AppVersion, 0)
	}
	a.Versions = append(a.Versions, version)
}

func (a *App) GetVersionById(id int64) *AppVersion {
	for _, v := range a.Versions {
		if id == 0 {
			return v
		}
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (a *App) DeleteVersion(version string) {
	for index, v := range a.Versions {
		if v.Version == version {
			a.Versions = append(a.Versions[:index], a.Versions[index+1:]...)
			return
		}
	}
}

func (a *AppUsecase) CheckCluster(_ context.Context) bool {
	_, err := GetKubeClientByInCluster()
	return err == nil
}

func (a *AppUsecase) GetAppConfigPath(appName, appVersionNumber string) string {
	appPath := utils.GetServerStoragePathByNames(AppPackage)
	return fmt.Sprintf("%s/%s-%s.yaml", appPath, appName, appVersionNumber)
}

func (a *AppUsecase) InstallBasicComponent(ctx context.Context, basicAppType BasicComponentAppType) ([]*App, []*AppRelease, error) {
	appPath := utils.GetServerStoragePathByNames(AppPackage)
	apps := make([]*App, 0)
	appReleases := make([]*AppRelease, 0)
	confApps := a.conf.Apps
	for _, v := range confApps {
		if strings.ToUpper(v.Type) != basicAppType.String() {
			continue
		}
		appchart := fmt.Sprintf("%s/%s-%s.tgz", appPath, v.Name, v.Version)
		if !utils.IsFileExist(appchart) {
			return nil, nil, errors.Errorf("appchart not found: %s", appchart)
		}
		app := &App{Name: v.Name}
		appVersion := &AppVersion{Chart: appchart, Version: v.Version}
		err := a.GetAppAndVersionInfo(ctx, app, appVersion)
		if err != nil {
			return nil, nil, err
		}
		app.AddVersion(appVersion)
		apps = append(apps, app)
		appConfigPath := a.GetAppConfigPath(app.Name, appVersion.Version)
		if !utils.IsFileExist(appConfigPath) {
			err = utils.WriteFile(appConfigPath, appVersion.DefaultConfig)
			if err != nil {
				return nil, nil, err
			}
		}
		appRelease := &AppRelease{
			AppId:       app.Id,
			VersionId:   appVersion.Id,
			ReleaseName: fmt.Sprintf("%s-%s", strings.ToLower(v.Type), v.Name),
			Namespace:   SystemNamespace,
			ConfigFile:  appConfigPath,
			Status:      AppReleaseSatus_BAppReleaseSatus_UNSPECIFIED,
		}
		appReleases = append(appReleases, appRelease)
	}
	return apps, appReleases, nil
}

func (a *AppUsecase) DeleteApp(ctx context.Context, app *App) error {
	appPath := utils.GetServerStoragePathByNames(AppPackage)
	err := os.Remove(appPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppUsecase) DeleteAppVersion(ctx context.Context, app *App, appVersion *AppVersion) error {
	appPath := utils.GetServerStoragePathByNames(AppPackage)
	for _, v := range app.Versions {
		if v.Chart == "" || appVersion.Id != v.Id {
			continue
		}
		chartPath := filepath.Join(appPath, v.Chart)
		if utils.IsFileExist(chartPath) {
			err := os.Remove(chartPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AppUsecase) GetAppAndVersionInfo(ctx context.Context, app *App, appVersion *AppVersion) error {
	charInfo, err := getLocalChartInfo(appVersion.Chart)
	if err != nil {
		return err
	}
	charInfoMetadata, err := json.Marshal(charInfo.Metadata)
	if err != nil {
		return err
	}
	app.Name = charInfo.Name
	app.Readme = charInfo.Readme
	app.Description = charInfo.Description
	app.Metadata = charInfoMetadata

	appVersion.DefaultConfig = charInfo.Config
	appVersion.Version = charInfo.Version
	appVersion.Chart = charInfo.Chart
	return nil
}

func (a *AppUsecase) AppRelease(ctx context.Context, app *App, appVersion *AppVersion, appRelease *AppRelease, appRepo *AppRepo) error {
	helmPkg, err := NewHelmPkg(a.log, appRelease.Namespace)
	if err != nil {
		return err
	}
	install, err := helmPkg.NewInstall()
	if err != nil {
		return err
	}
	if appRepo != nil && appVersion.Chart == "" {
		appPath := utils.GetServerStoragePathByNames(AppPackage)
		pullClient, err := helmPkg.NewPull()
		if err != nil {
			return err
		}
		err = helmPkg.RunPull(pullClient, appPath, fmt.Sprintf("%s/%s", appRepo.Name, app.Name))
		if err != nil {
			return err
		}
		appVersion.Chart = filepath.Join(appPath, fmt.Sprintf("%s-%s.tgz", app.Name, appVersion.Version))
	}
	install.ReleaseName = appRelease.ReleaseName
	install.Namespace = appRelease.Namespace
	install.CreateNamespace = true
	install.GenerateName = true
	install.Version = appVersion.Version
	install.DryRun = appRelease.Dryrun
	install.Atomic = appRelease.Atomic
	install.Wait = appRelease.Wait
	release, err := helmPkg.RunInstall(ctx, install, appVersion.Chart, &helmValues.Options{
		ValueFiles: []string{appRelease.ConfigFile},
		Values:     []string{appRelease.Config},
	})
	appRelease.Logs = helmPkg.GetLogs()
	if err != nil {
		return err
	}
	if release != nil {
		appRelease.ReleaseName = release.Name
		if release.Info != nil {
			appRelease.Notes = release.Info.Notes
		}
		return nil
	}
	return nil
}

func (a *AppUsecase) GetAppReleaseResources(ctx context.Context, appRelease *AppRelease) ([]*AppReleaseResource, error) {
	helmPkg, err := NewHelmPkg(a.log, appRelease.Namespace)
	if err != nil {
		return nil, err
	}
	releaseInfo, err := helmPkg.GetReleaseInfo(appRelease.ReleaseName)
	if err != nil {
		return nil, err
	}
	if releaseInfo == nil || releaseInfo.Manifest == "" {
		return nil, errors.New("release not found")
	}
	clusterClient, err := GetKubeClientByKubeConfig()
	if err != nil {
		return nil, err
	}
	events, err := clusterClient.CoreV1().Events(appRelease.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	resources := strings.Split(releaseInfo.Manifest, "---")
	appReleaseResources := make([]*AppReleaseResource, 0)
	for _, resource := range resources {
		if resource == "" {
			continue
		}
		obj := unstructured.Unstructured{}
		err = yaml.Unmarshal([]byte(resource), &obj.Object)
		if err != nil {
			return nil, err
		}
		lableStr := ""
		for k, v := range obj.GetLabels() {
			lableStr += fmt.Sprintf("%s=%s,", k, v)
		}
		lableStr = strings.TrimRight(lableStr, ",")
		appReleaseResource := &AppReleaseResource{
			Id:        uuid.NewString(),
			ReleaseId: appRelease.Id,
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
			Kind:      obj.GetKind(),
			Lables:    lableStr,
			Manifest:  resource,
			Status:    AppReleaseResourceStatus_SUCCESSFUL,
		}
		eventsStr := ""
		for _, event := range events.Items {
			if event.InvolvedObject.Name == obj.GetName() {
				if event.Type != "Normal" {
					appReleaseResource.Status = AppReleaseResourceStatus_UNHEALTHY
				}
				eventsStr += fmt.Sprintf("- Type: %s, Reason: %s, Message: %s\n", event.Type, event.Reason, event.Message)
			}
		}
		appReleaseResource.Events = eventsStr
		appReleaseResources = append(appReleaseResources, appReleaseResource)

		podLableSelector := ""
		if strings.EqualFold(appReleaseResource.Kind, Deployment) {
			deploymentResource, err := clusterClient.AppsV1().Deployments(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			for k, v := range deploymentResource.Spec.Selector.MatchLabels {
				podLableSelector += fmt.Sprintf("%s=%s,", k, v)
			}
			podLableSelector = strings.TrimRight(podLableSelector, ",")
		}
		if strings.EqualFold(appReleaseResource.Kind, StatefulSet) {
			statefulSetResource, err := clusterClient.AppsV1().StatefulSets(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			for k, v := range statefulSetResource.Spec.Selector.MatchLabels {
				podLableSelector += fmt.Sprintf("%s=%s,", k, v)
			}
			podLableSelector = strings.TrimRight(podLableSelector, ",")
		}
		if strings.EqualFold(appReleaseResource.Kind, DaemonSet) {
			daemonSetResource, err := clusterClient.AppsV1().DaemonSets(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			for k, v := range daemonSetResource.Spec.Selector.MatchLabels {
				podLableSelector += fmt.Sprintf("%s=%s,", k, v)
			}
			podLableSelector = strings.TrimRight(podLableSelector, ",")
		}
		if podLableSelector != "" {
			podResources, err := clusterClient.CoreV1().Pods(appReleaseResource.Namespace).List(ctx, metav1.ListOptions{
				LabelSelector: podLableSelector,
			})
			if err != nil {
				return nil, err
			}
			for _, pod := range podResources.Items {
				lableStr := ""
				for k, v := range pod.Labels {
					lableStr += fmt.Sprintf("%s=%s,", k, v)
				}
				lableStr = strings.TrimRight(lableStr, ",")
				appReleaseResource := &AppReleaseResource{
					Id:        uuid.NewString(),
					ReleaseId: appRelease.Id,
					Name:      fmt.Sprintf("%s-%s-%s", obj.GetKind(), obj.GetName(), pod.Name),
					Namespace: pod.Namespace,
					Kind:      Pod,
					Lables:    lableStr,
					Status:    AppReleaseResourceStatus_SUCCESSFUL,
				}
				eventsStr := ""
				for _, event := range events.Items {
					if string(event.InvolvedObject.UID) == string(pod.UID) {
						if event.Type != coreV1.EventTypeNormal {
							appReleaseResource.Status = AppReleaseResourceStatus_UNHEALTHY
						}
						eventsStr += fmt.Sprintf("- Type: %s, Reason: %s, Message: %s\n", event.Type, event.Reason, event.Message)
					}
				}
				if pod.Status.Phase != coreV1.PodRunning {
					appReleaseResource.Status = AppReleaseResourceStatus_UNHEALTHY
				}
				for _, containerStatus := range pod.Status.ContainerStatuses {
					if !containerStatus.Ready {
						appReleaseResource.Status = AppReleaseResourceStatus_UNHEALTHY
						break
					}
				}
				appReleaseResource.Events = eventsStr
				appReleaseResources = append(appReleaseResources, appReleaseResource)
			}
		}
	}
	return appReleaseResources, nil
}

func (a *AppUsecase) ReloadAppReleaseResource(ctx context.Context, appReleaseResource *AppReleaseResource) error {
	clusterClient, err := GetKubeClientByKubeConfig()
	if err != nil {
		return err
	}
	if strings.EqualFold(appReleaseResource.Kind, Pod) {
		err = clusterClient.CoreV1().Pods(appReleaseResource.GetNamespace()).Delete(ctx, appReleaseResource.GetName(), metav1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	podLableSelector := ""
	if strings.EqualFold(appReleaseResource.Kind, Deployment) {
		deploymentResource, err := clusterClient.AppsV1().Deployments(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}
		for k, v := range deploymentResource.Spec.Selector.MatchLabels {
			podLableSelector += fmt.Sprintf("%s=%s,", k, v)
		}
		podLableSelector = strings.TrimRight(podLableSelector, ",")
	}
	if strings.EqualFold(appReleaseResource.Kind, StatefulSet) {
		statefulSetResource, err := clusterClient.AppsV1().StatefulSets(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}
		for k, v := range statefulSetResource.Spec.Selector.MatchLabels {
			podLableSelector += fmt.Sprintf("%s=%s,", k, v)
		}
		podLableSelector = strings.TrimRight(podLableSelector, ",")
	}
	if strings.EqualFold(appReleaseResource.Kind, DaemonSet) {
		daemonSetResource, err := clusterClient.AppsV1().DaemonSets(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}
		for k, v := range daemonSetResource.Spec.Selector.MatchLabels {
			podLableSelector += fmt.Sprintf("%s=%s,", k, v)
		}
		podLableSelector = strings.TrimRight(podLableSelector, ",")
	}
	if podLableSelector == "" {
		return nil
	}
	podResources, err := clusterClient.CoreV1().Pods(appReleaseResource.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: podLableSelector,
	})
	if err != nil {
		return err
	}
	for _, pod := range podResources.Items {
		err = clusterClient.CoreV1().Pods(appReleaseResource.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AppUsecase) DeleteAppRelease(ctx context.Context, appRelease *AppRelease) error {
	helmPkg, err := NewHelmPkg(a.log, appRelease.Namespace)
	if err != nil {
		return err
	}
	uninstall, err := helmPkg.NewUninstall()
	if err != nil {
		return err
	}
	uninstall.KeepHistory = false
	uninstall.DryRun = appRelease.Dryrun
	uninstall.Wait = appRelease.Wait
	resp, err := helmPkg.RunUninstall(uninstall, appRelease.ReleaseName)
	appRelease.Logs = helmPkg.GetLogs()
	if err != nil {
		return errors.WithMessage(err, "uninstall fail")
	}
	if resp != nil && resp.Release != nil && resp.Release.Info != nil {
		appRelease.Status = AppReleaseSatus_RUNNING
	}
	appRelease.Notes = resp.Info
	return nil
}

func (a *AppUsecase) AddAppRepo(ctx context.Context, repo *AppRepo) error {
	settings := cli.New()
	res, err := helmrepo.NewChartRepository(&helmrepo.Entry{
		Name: repo.Name,
		URL:  repo.Url,
	}, getter.All(settings))
	if err != nil {
		return err
	}
	res.CachePath = utils.GetServerStoragePathByNames(AppPackage)
	indexFile, err := res.DownloadIndexFile()
	if err != nil {
		return err
	}
	repo.IndexPath = indexFile
	return nil
}

func (a *AppUsecase) GetAppsByRepo(ctx context.Context, repo *AppRepo) ([]*App, error) {
	index, err := helmrepo.LoadIndexFile(repo.IndexPath)
	if err != nil {
		return nil, err
	}
	apps := make([]*App, 0)
	for chartName, chartVersions := range index.Entries {
		app := &App{Name: chartName, AppRepoId: repo.Id, Versions: make([]*AppVersion, 0)}
		for _, chartMatedata := range chartVersions {
			if len(chartMatedata.URLs) == 0 {
				return nil, errors.New("chart urls is empty")
			}
			app.Icon = chartMatedata.Icon
			app.Description = chartMatedata.Description
			appVersion := &AppVersion{Name: chartMatedata.Name, Chart: chartMatedata.URLs[0], Version: chartMatedata.Version}
			app.Versions = append(app.Versions, appVersion)
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (a *AppUsecase) GetAppDetailByRepo(ctx context.Context, apprepo *AppRepo, appName, version string) (*App, error) {
	index, err := helmrepo.LoadIndexFile(apprepo.IndexPath)
	if err != nil {
		return nil, err
	}
	app := &App{Name: appName, AppRepoId: apprepo.Id, Versions: make([]*AppVersion, 0)}
	for chartName, chartVersions := range index.Entries {
		if chartName != appName {
			continue
		}
		for i, chartMatedata := range chartVersions {
			if len(chartMatedata.URLs) == 0 {
				return nil, errors.New("chart urls is empty")
			}
			app.Icon = chartMatedata.Icon
			app.Name = chartName
			app.Description = chartMatedata.Description
			appVersion := &AppVersion{Name: chartMatedata.Name, Chart: chartMatedata.URLs[0], Version: chartMatedata.Version}
			if (version == "" && i == 0) || (version != "" && version == chartMatedata.Version) {
				err = a.GetAppAndVersionInfo(ctx, app, appVersion)
				if err != nil {
					return nil, err
				}
			}
			app.Versions = append(app.Versions, appVersion)
		}
	}
	return app, nil
}
