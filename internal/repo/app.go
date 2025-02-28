package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

type AppRepo struct {
	log *log.Helper
}

func NewAppRepo(logger log.Logger) biz.AppRepoInterface {
	return &AppRepo{
		log: log.NewHelper(logger),
	}
}

func (r *AppRepo) CheckCluster(ctx context.Context) bool {
	_, err := GetKubeClientByKubeConfig()
	return err == nil
}

func (r *AppRepo) GetAppAndVersionInfo(ctx context.Context, app *biz.App, appVersion *biz.AppVersion) error {
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

func (r *AppRepo) AppRelease(ctx context.Context, app *biz.App, appVersion *biz.AppVersion, appRelease *biz.AppRelease, appRepo *biz.AppRepo) error {
	helmPkg, err := NewHelmPkg(r.log, appRelease.Namespace)
	if err != nil {
		return err
	}
	install, err := helmPkg.NewInstall()
	if err != nil {
		return err
	}
	if appRepo != nil && appVersion.Chart == "" {
		appPath := utils.GetServerStoragePathByNames(biz.AppPackage)
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

func (r *AppRepo) GetAppReleaseResources(ctx context.Context, appRelease *biz.AppRelease) ([]*biz.AppReleaseResource, error) {
	helmPkg, err := NewHelmPkg(r.log, appRelease.Namespace)
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
	appReleaseResources := make([]*biz.AppReleaseResource, 0)
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
		appReleaseResource := &biz.AppReleaseResource{
			Id:        uuid.NewString(),
			ReleaseId: appRelease.Id,
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
			Kind:      obj.GetKind(),
			Lables:    lableStr,
			Manifest:  resource,
			Status:    biz.AppReleaseResourceStatus_SUCCESSFUL,
		}
		eventsStr := ""
		for _, event := range events.Items {
			if event.InvolvedObject.Name == obj.GetName() {
				if event.Type != "Normal" {
					appReleaseResource.Status = biz.AppReleaseResourceStatus_UNHEALTHY
				}
				eventsStr += fmt.Sprintf("- Type: %s, Reason: %s, Message: %s\n", event.Type, event.Reason, event.Message)
			}
		}
		appReleaseResource.Events = eventsStr
		appReleaseResources = append(appReleaseResources, appReleaseResource)

		podLableSelector := ""
		if strings.EqualFold(appReleaseResource.Kind, biz.Deployment) {
			deploymentResource, err := clusterClient.AppsV1().Deployments(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			for k, v := range deploymentResource.Spec.Selector.MatchLabels {
				podLableSelector += fmt.Sprintf("%s=%s,", k, v)
			}
			podLableSelector = strings.TrimRight(podLableSelector, ",")
		}
		if strings.EqualFold(appReleaseResource.Kind, biz.StatefulSet) {
			statefulSetResource, err := clusterClient.AppsV1().StatefulSets(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			for k, v := range statefulSetResource.Spec.Selector.MatchLabels {
				podLableSelector += fmt.Sprintf("%s=%s,", k, v)
			}
			podLableSelector = strings.TrimRight(podLableSelector, ",")
		}
		if strings.EqualFold(appReleaseResource.Kind, biz.DaemonSet) {
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
				appReleaseResource := &biz.AppReleaseResource{
					Id:        uuid.NewString(),
					ReleaseId: appRelease.Id,
					Name:      fmt.Sprintf("%s-%s-%s", obj.GetKind(), obj.GetName(), pod.Name),
					Namespace: pod.Namespace,
					Kind:      biz.Pod,
					Lables:    lableStr,
					Status:    biz.AppReleaseResourceStatus_SUCCESSFUL,
				}
				eventsStr := ""
				for _, event := range events.Items {
					if string(event.InvolvedObject.UID) == string(pod.UID) {
						if event.Type != coreV1.EventTypeNormal {
							appReleaseResource.Status = biz.AppReleaseResourceStatus_UNHEALTHY
						}
						eventsStr += fmt.Sprintf("- Type: %s, Reason: %s, Message: %s\n", event.Type, event.Reason, event.Message)
					}
				}
				if pod.Status.Phase != coreV1.PodRunning {
					appReleaseResource.Status = biz.AppReleaseResourceStatus_UNHEALTHY
				}
				for _, containerStatus := range pod.Status.ContainerStatuses {
					if !containerStatus.Ready {
						appReleaseResource.Status = biz.AppReleaseResourceStatus_UNHEALTHY
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

func (r *AppRepo) ReloadAppReleaseResource(ctx context.Context, appReleaseResource *biz.AppReleaseResource) error {
	clusterClient, err := GetKubeClientByKubeConfig()
	if err != nil {
		return err
	}
	if strings.EqualFold(appReleaseResource.Kind, biz.Pod) {
		err = clusterClient.CoreV1().Pods(appReleaseResource.GetNamespace()).Delete(ctx, appReleaseResource.GetName(), metav1.DeleteOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	podLableSelector := ""
	if strings.EqualFold(appReleaseResource.Kind, biz.Deployment) {
		deploymentResource, err := clusterClient.AppsV1().Deployments(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}
		for k, v := range deploymentResource.Spec.Selector.MatchLabels {
			podLableSelector += fmt.Sprintf("%s=%s,", k, v)
		}
		podLableSelector = strings.TrimRight(podLableSelector, ",")
	}
	if strings.EqualFold(appReleaseResource.Kind, biz.StatefulSet) {
		statefulSetResource, err := clusterClient.AppsV1().StatefulSets(appReleaseResource.Namespace).Get(ctx, appReleaseResource.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}
		for k, v := range statefulSetResource.Spec.Selector.MatchLabels {
			podLableSelector += fmt.Sprintf("%s=%s,", k, v)
		}
		podLableSelector = strings.TrimRight(podLableSelector, ",")
	}
	if strings.EqualFold(appReleaseResource.Kind, biz.DaemonSet) {
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

func (r *AppRepo) DeleteAppRelease(ctx context.Context, appRelease *biz.AppRelease) error {
	helmPkg, err := NewHelmPkg(r.log, appRelease.Namespace)
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
		appRelease.Status = biz.AppReleaseSatus_RUNNING
	}
	appRelease.Notes = resp.Info
	return nil
}

func (r *AppRepo) AddAppRepo(ctx context.Context, repo *biz.AppRepo) error {
	settings := cli.New()
	res, err := helmrepo.NewChartRepository(&helmrepo.Entry{
		Name: repo.Name,
		URL:  repo.Url,
	}, getter.All(settings))
	if err != nil {
		return err
	}
	res.CachePath = utils.GetServerStoragePathByNames(biz.AppPackage)
	indexFile, err := res.DownloadIndexFile()
	if err != nil {
		return err
	}
	repo.IndexPath = indexFile
	return nil
}

func (r *AppRepo) GetAppsByRepo(ctx context.Context, apprepo *biz.AppRepo) ([]*biz.App, error) {
	index, err := helmrepo.LoadIndexFile(apprepo.IndexPath)
	if err != nil {
		return nil, err
	}
	apps := make([]*biz.App, 0)
	for chartName, chartVersions := range index.Entries {
		app := &biz.App{Name: chartName, AppRepoId: apprepo.Id, Versions: make([]*biz.AppVersion, 0)}
		for _, chartMatedata := range chartVersions {
			if len(chartMatedata.URLs) == 0 {
				return nil, errors.New("chart urls is empty")
			}
			app.Icon = chartMatedata.Icon
			app.Description = chartMatedata.Description
			appVersion := &biz.AppVersion{Name: chartMatedata.Name, Chart: chartMatedata.URLs[0], Version: chartMatedata.Version}
			app.Versions = append(app.Versions, appVersion)
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *AppRepo) GetAppDetailByRepo(ctx context.Context, apprepo *biz.AppRepo, appName, version string) (*biz.App, error) {
	index, err := helmrepo.LoadIndexFile(apprepo.IndexPath)
	if err != nil {
		return nil, err
	}
	app := &biz.App{Name: appName, AppRepoId: apprepo.Id, Versions: make([]*biz.AppVersion, 0)}
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
			appVersion := &biz.AppVersion{Name: chartMatedata.Name, Chart: chartMatedata.URLs[0], Version: chartMatedata.Version}
			if (version == "" && i == 0) || (version != "" && version == chartMatedata.Version) {
				err = r.GetAppAndVersionInfo(ctx, app, appVersion)
				if err != nil {
					return nil, err
				}
			}
			app.Versions = append(app.Versions, appVersion)
		}
	}
	return app, nil
}

func GetKubeClientByKubeConfig(KubeConfigPaths ...string) (clientset *kubernetes.Clientset, err error) {
	var KubeConfigPath string
	if len(KubeConfigPaths) == 0 {
		KubeConfigPath = clientcmd.RecommendedHomeFile
	} else {
		KubeConfigPath = KubeConfigPaths[0]
	}
	config, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes by kubeconfig client failed")
	}
	return client, nil
}
