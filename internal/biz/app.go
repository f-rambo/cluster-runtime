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
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	helmrepo "helm.sh/helm/v3/pkg/repo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AppPackage string = "app"

type AppUsecase struct {
	log  *log.Helper
	conf *conf.Bootstrap
}

func NewAppUseCase(logger log.Logger) *AppUsecase {
	a := &AppUsecase{
		log: log.NewHelper(logger),
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

func (a *AppUsecase) InstallBasicComponent(ctx context.Context, basicAppType BasicComponentAppType) ([]*App, []*AppRelease, error) {
	appPath := utils.GetServerStoragePathByNames(AppPackage)
	configPath := utils.GetFromContextByKey(ctx, utils.ConfDirKey)
	apps := make([]*App, 0)
	appReleases := make([]*AppRelease, 0)
	confApps := a.conf.Apps
	for _, v := range confApps {
		if strings.ToUpper(v.Type) != basicAppType.String() {
			continue
		}
		appchart := fmt.Sprintf("%s/%s-%s.tgz", appPath, v.Name, v.Version)
		if !utils.IsFileExist(appchart) {
			return nil, nil, fmt.Errorf("appchart not found: %s", appchart)
		}
		app := &App{Name: v.Name}
		appVersion := &AppVersion{Chart: appchart, Version: v.Version}
		err := a.GetAppAndVersionInfo(ctx, app, appVersion)
		if err != nil {
			return nil, nil, err
		}
		app.AddVersion(appVersion)
		apps = append(apps, app)
		appConfigPath := fmt.Sprintf("%s/%s-%s.yaml", configPath, v.Name, v.Version)
		if utils.IsFileExist(appConfigPath) {
			appConfig, err := os.ReadFile(appConfigPath)
			if err != nil {
				return nil, nil, err
			}
			appVersion.DefaultConfig = string(appConfig)
		}
		appRelease := &AppRelease{
			ReleaseName: fmt.Sprintf("%s-%s", v.Name, v.Version),
			AppId:       app.Id,
			VersionId:   appVersion.Id,
			Namespace:   a.conf.Server.Namespace,
			Config:      appVersion.DefaultConfig,
			Status:      AppReleaseSatus_PENDING,
			Wait:        true,
		}
		appReleases = append(appReleases, appRelease)
	}
	return apps, appReleases, nil
}

func (a *AppUsecase) GetPodResources(ctx context.Context, appRelease *AppRelease) ([]*AppReleaseResource, error) {
	resources := make([]*AppReleaseResource, 0)
	clusterClient, err := GetKubeClientByInCluster()
	if err != nil {
		return nil, err
	}
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", appRelease.ReleaseName)
	podResources, _ := clusterClient.CoreV1().Pods(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if podResources != nil && len(podResources.Items) > 0 {
		for _, pod := range podResources.Items {
			resource := &AppReleaseResource{
				Name:      pod.Name,
				Kind:      "Pod",
				StartedAt: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status:    []string{string(pod.Status.Phase)},
				Events:    make([]string, 0),
			}
			events, _ := clusterClient.CoreV1().Events(appRelease.Namespace).List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name),
			})
			if events != nil && len(events.Items) > 0 {
				for _, event := range events.Items {
					resource.Events = append(resource.Events, event.Message)
				}
			}
			resources = append(resources, resource)
		}
	}
	return resources, nil
}

func (a *AppUsecase) GetNetResouces(ctx context.Context, appRelease *AppRelease) ([]*AppReleaseResource, error) {
	resources := make([]*AppReleaseResource, 0)
	clusterClient, err := GetKubeClientByInCluster()
	if err != nil {
		return nil, err
	}
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", appRelease.ReleaseName)
	serviceResources, _ := clusterClient.CoreV1().Services(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if serviceResources != nil && len(serviceResources.Items) > 0 {
		for _, service := range serviceResources.Items {
			port := ""
			for _, v := range service.Spec.Ports {
				port = fmt.Sprintf("%s %d:%d/%s", port, v.Port, v.NodePort, v.Protocol)
			}
			externalIPs := ""
			for _, v := range service.Spec.ExternalIPs {
				externalIPs = fmt.Sprintf("%s,%s", externalIPs, v)
			}
			resource := &AppReleaseResource{
				Name:      service.Name,
				Kind:      "Service",
				StartedAt: service.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"Type: " + string(service.Spec.Type),
					"ClusterIP: " + service.Spec.ClusterIP,
					"ExternalIP: " + externalIPs,
					"Port: " + port,
				},
			}
			resources = append(resources, resource)
		}
	}
	ingressResources, _ := clusterClient.NetworkingV1beta1().Ingresses(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if ingressResources != nil && len(ingressResources.Items) > 0 {
		for _, ingress := range ingressResources.Items {
			class := ""
			if ingress.Spec.IngressClassName != nil {
				class = *ingress.Spec.IngressClassName
			}
			hosts := ""
			for _, v := range ingress.Spec.Rules {
				hosts = fmt.Sprintf("%s,%s", hosts, v.Host)
			}
			ports := ""
			for _, v := range ingress.Spec.TLS {
				for _, v := range v.Hosts {
					ports = fmt.Sprintf("%s,%s", ports, v)
				}
			}
			loadBalancerIP := ""
			for _, v := range ingress.Status.LoadBalancer.Ingress {
				loadBalancerIP = fmt.Sprintf("%s,%s", loadBalancerIP, v.IP)
			}
			resource := &AppReleaseResource{
				Name:      ingress.Name,
				Kind:      "Ingress",
				StartedAt: ingress.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"class: " + class,
					"hosts: " + hosts,
					"address: " + loadBalancerIP,
					"ports: " + ports,
				},
			}
			resources = append(resources, resource)
		}
	}
	return resources, nil
}

func (a *AppUsecase) GetAppsReouces(ctx context.Context, appRelease *AppRelease) ([]*AppReleaseResource, error) {
	resources := make([]*AppReleaseResource, 0)
	clusterClient, err := GetKubeClientByInCluster()
	if err != nil {
		return nil, err
	}
	labelSelector := fmt.Sprintf("app.kubernetes.io/instance=%s", appRelease.ReleaseName)
	deploymentResources, _ := clusterClient.AppsV1().Deployments(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if deploymentResources != nil && len(deploymentResources.Items) > 0 {
		for _, deployment := range deploymentResources.Items {
			resource := &AppReleaseResource{
				Name:      deployment.Name,
				Kind:      "Deployment",
				StartedAt: deployment.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"Ready: " + cast.ToString(deployment.Status.ReadyReplicas),
					"Up-to-date: " + cast.ToString(deployment.Status.UpdatedReplicas),
					"Available: " + cast.ToString(deployment.Status.AvailableReplicas),
				},
			}
			resources = append(resources, resource)
		}
	}

	statefulSetResources, _ := clusterClient.AppsV1().StatefulSets(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if statefulSetResources != nil && len(statefulSetResources.Items) > 0 {
		for _, statefulSet := range statefulSetResources.Items {
			resource := &AppReleaseResource{
				Name:      statefulSet.Name,
				Kind:      "StatefulSet",
				StartedAt: statefulSet.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"Ready: " + cast.ToString(statefulSet.Status.ReadyReplicas),
					"Up-to-date: " + cast.ToString(statefulSet.Status.UpdatedReplicas),
					"Available: " + cast.ToString(statefulSet.Status.AvailableReplicas),
				},
			}
			resources = append(resources, resource)
		}
	}

	deamonsetResources, _ := clusterClient.AppsV1().DaemonSets(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if deamonsetResources != nil && len(deamonsetResources.Items) > 0 {
		for _, deamonset := range deamonsetResources.Items {
			resource := &AppReleaseResource{
				Name:      deamonset.Name,
				Kind:      "Deamonset",
				StartedAt: deamonset.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"Desired: " + cast.ToString(deamonset.Status.DesiredNumberScheduled),
					"Current: " + cast.ToString(deamonset.Status.CurrentNumberScheduled),
					"Ready: " + cast.ToString(deamonset.Status.NumberReady),
					"Up-to-date: " + cast.ToString(deamonset.Status.UpdatedNumberScheduled),
					"Available: " + cast.ToString(deamonset.Status.NumberAvailable),
				},
			}
			resources = append(resources, resource)
		}
	}

	jobResources, _ := clusterClient.BatchV1().Jobs(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if jobResources != nil && len(jobResources.Items) > 0 {
		for _, job := range jobResources.Items {
			resource := &AppReleaseResource{
				Name:      job.Name,
				Kind:      "Job",
				StartedAt: job.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"Completions: " + cast.ToString(job.Spec.Completions),
					"Parallelism: " + cast.ToString(job.Spec.Parallelism),
					"BackoffLimit: " + cast.ToString(job.Spec.BackoffLimit),
				},
			}
			resources = append(resources, resource)
		}
	}

	cronjobResources, _ := clusterClient.BatchV1beta1().CronJobs(appRelease.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if cronjobResources != nil && len(cronjobResources.Items) > 0 {
		for _, cronjob := range cronjobResources.Items {
			suspend := *cronjob.Spec.Suspend // Dereference the pointer to bool
			resource := &AppReleaseResource{
				Name:      cronjob.Name,
				Kind:      "Cronjob",
				StartedAt: cronjob.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Status: []string{
					"Schedule: " + cronjob.Spec.Schedule,
					"Suspend: " + cast.ToString(suspend),
					"Active: " + cast.ToString(len(cronjob.Status.Active)),
					"Last Schedule: " + cronjob.Status.LastScheduleTime.Format("2006-01-02 15:04:05"),
				},
			}
			resources = append(resources, resource)
		}
	}
	return resources, nil
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
	release, err := helmPkg.RunInstall(ctx, install, appVersion.Chart, appRelease.Config)
	appRelease.Logs = helmPkg.GetLogs()
	if err != nil {
		return err
	}
	if release != nil {
		appRelease.ReleaseName = release.Name
		appRelease.Manifest = strings.TrimSpace(release.Manifest)
		if release.Info != nil {
			appRelease.Status = AppReleaseSatus_RUNNING
			appRelease.Notes = release.Info.Notes
		}
		return nil
	}
	appRelease.Status = AppReleaseSatus_FAILED
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
