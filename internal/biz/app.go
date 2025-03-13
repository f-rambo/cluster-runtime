package biz

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/component"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/conf"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
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

type AppRepoInterface interface {
	CheckCluster(context.Context) bool
	GetAppAndVersionInfo(context.Context, *App, *AppVersion) error
	AppRelease(context.Context, *App, *AppVersion, *AppRelease, *AppRepo) error
	GetAppReleaseResources(context.Context, *AppRelease) ([]*AppReleaseResource, error)
	ReloadAppReleaseResource(context.Context, *AppReleaseResource) error
	DeleteAppRelease(context.Context, *AppRelease) error
	AddAppRepo(context.Context, *AppRepo) error
	GetAppsByRepo(context.Context, *AppRepo) ([]*App, error)
	GetAppDetailByRepo(ctx context.Context, apprepo *AppRepo, appName, version string) (*App, error)
}

type AppUsecase struct {
	repo AppRepoInterface
	log  *log.Helper
	conf *conf.Bootstrap
}

func NewAppUseCase(conf *conf.Bootstrap, repo AppRepoInterface, logger log.Logger) *AppUsecase {
	a := &AppUsecase{
		conf: conf,
		repo: repo,
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

func (a *AppUsecase) CheckCluster(ctx context.Context) bool {
	return a.repo.CheckCluster(ctx)
}

func (a *AppUsecase) GetAppConfigPath(appName, appVersionNumber string) string {
	appPath := utils.GetServerStoragePathByNames(AppPackage)
	return fmt.Sprintf("%s/%s-%s.yaml", appPath, appName, appVersionNumber)
}

func (a *AppUsecase) InstallBasicComponent(ctx context.Context, cluster *Cluster, basicAppType BasicComponentAppType) ([]*App, []*AppRelease, error) {
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
		appConfig := fmt.Sprintf("%s/%s-%s.yaml", appPath, v.Name, v.Version)
		err := component.TransferredMeaning(cluster, appConfig)
		if err != nil {
			return nil, nil, err
		}
		app := &App{Name: v.Name}
		appVersion := &AppVersion{Chart: appchart, Version: v.Version}
		err = a.GetAppAndVersionInfo(ctx, app, appVersion)
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

func (a *AppUsecase) SettingBasicComponent(ctx context.Context, cluster *Cluster) error {
	// Read compoment config and update config apply to cluster. example: ceph-cluster... by nodes info
	return nil
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
	return a.repo.GetAppAndVersionInfo(ctx, app, appVersion)
}

func (a *AppUsecase) AppRelease(ctx context.Context, app *App, appVersion *AppVersion, appRelease *AppRelease, appRepo *AppRepo) error {
	return a.repo.AppRelease(ctx, app, appVersion, appRelease, appRepo)
}

func (a *AppUsecase) GetAppReleaseResources(ctx context.Context, appRelease *AppRelease) ([]*AppReleaseResource, error) {
	return a.repo.GetAppReleaseResources(ctx, appRelease)
}

func (a *AppUsecase) ReloadAppReleaseResource(ctx context.Context, appReleaseResource *AppReleaseResource) error {
	return a.repo.ReloadAppReleaseResource(ctx, appReleaseResource)
}

func (a *AppUsecase) DeleteAppRelease(ctx context.Context, appRelease *AppRelease) error {
	return a.repo.DeleteAppRelease(ctx, appRelease)
}

func (a *AppUsecase) AddAppRepo(ctx context.Context, repo *AppRepo) error {
	return a.repo.AddAppRepo(ctx, repo)
}

func (a *AppUsecase) GetAppsByRepo(ctx context.Context, repo *AppRepo) ([]*App, error) {
	return a.repo.GetAppsByRepo(ctx, repo)
}

func (a *AppUsecase) GetAppDetailByRepo(ctx context.Context, apprepo *AppRepo, appName, version string) (*App, error) {
	return a.repo.GetAppDetailByRepo(ctx, apprepo, appName, version)
}
