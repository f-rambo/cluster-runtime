package interfaces

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/app"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/api/common"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AppInterface struct {
	appApi.UnimplementedAppInterfaceServer
	appUC *biz.AppUsecase
	log   *log.Helper
}

func NewAppInterface(uc *biz.AppUsecase, logger log.Logger) *AppInterface {
	return &AppInterface{
		appUC: uc,
		log:   log.NewHelper(logger),
	}
}

func (a *AppInterface) UploadApp(ctx context.Context, req *appApi.FileUploadRequest) (*appApi.GetAppAndVersionInfo, error) {
	var fileExt string = ".tgz"
	if filepath.Ext(req.GetFileName()) != fileExt {
		return nil, errors.New("file type is not supported")
	}
	appPath := utils.GetServerStoragePathByNames(biz.AppPackage)
	fileName, err := a.upload(appPath, req.GetFileName(), req.GetChunk())
	if err != nil {
		return nil, err
	}
	appTmpChartPath := fmt.Sprintf("%s/%s", appPath, fileName)
	app := &biz.App{Versions: make([]*biz.AppVersion, 0)}
	appVersion := &biz.AppVersion{Chart: appTmpChartPath}
	err = a.appUC.GetAppAndVersionInfo(ctx, app, appVersion)
	if err != nil {
		return nil, err
	}
	app.AddVersion(appVersion)
	appChartPath := fmt.Sprintf("%s/%s-%s%s", appPath, app.Name, appVersion.Version, fileExt)
	if utils.IsFileExist(appChartPath) {
		err = os.Remove(appChartPath)
		if err != nil {
			return nil, err
		}
	}
	err = os.Rename(appTmpChartPath, appChartPath)
	if err != nil {
		return nil, err
	}
	return &appApi.GetAppAndVersionInfo{App: app}, nil
}

func (a *AppInterface) upload(path, filename, chunk string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(chunk[strings.IndexByte(chunk, ',')+1:])
	if err != nil {
		return "", err
	}
	file, err := utils.NewFile(path, filename, false)
	if err != nil {
		return "", err
	}
	defer func() {
		if file == nil {
			return
		}
		err := file.Close()
		if err != nil {
			a.log.Error(err)
		}
	}()
	err = file.Write(data)
	if err != nil {
		return "", err
	}
	return file.GetFileName(), nil
}

func (a *AppInterface) CheckCluster(ctx context.Context, _ *emptypb.Empty) (*appApi.CheckClusterResponse, error) {
	ok := a.appUC.CheckCluster(ctx)
	return &appApi.CheckClusterResponse{Ok: ok}, nil
}

func (a *AppInterface) InstallBasicComponent(ctx context.Context, param *appApi.InstallBasicComponentReq) (*appApi.InstallBasicComponentResponse, error) {
	if param.BasicComponentAppType == 0 {
		return nil, errors.New("basic component app type is empty")
	}
	apps, appReleases, err := a.appUC.InstallBasicComponent(ctx, param.Cluster, param.BasicComponentAppType)
	if err != nil {
		return nil, err
	}
	appItems := &appApi.InstallBasicComponentResponse{Apps: apps, Releases: appReleases}
	return appItems, nil
}

func (a *AppInterface) GetAppReleaseResources(ctx context.Context, appRelease *biz.AppRelease) (*appApi.AppReleaseResourceItems, error) {
	appReleaseResources, err := a.appUC.GetAppReleaseResources(ctx, appRelease)
	if err != nil {
		return nil, err
	}
	return &appApi.AppReleaseResourceItems{Resources: appReleaseResources}, nil
}

func (a *AppInterface) DeleteApp(ctx context.Context, app *biz.App) (*common.Msg, error) {
	err := a.appUC.DeleteApp(ctx, app)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (a *AppInterface) DeleteAppVersion(ctx context.Context, deleteAppVersionReq *appApi.DeleteAppVersionReq) (*common.Msg, error) {
	err := a.appUC.DeleteAppVersion(ctx, deleteAppVersionReq.App, deleteAppVersionReq.Version)
	return common.Response(), err
}

func (a *AppInterface) GetAppAndVersionInfo(ctx context.Context, appAndVersionInfo *appApi.GetAppAndVersionInfo) (*appApi.GetAppAndVersionInfo, error) {
	err := a.appUC.GetAppAndVersionInfo(ctx, appAndVersionInfo.App, appAndVersionInfo.Version)
	return appAndVersionInfo, err
}

func (a *AppInterface) AppRelease(ctx context.Context, appRelease *appApi.AppReleaseReq) (*biz.AppRelease, error) {
	err := a.appUC.AppRelease(ctx, appRelease.App, appRelease.Version, appRelease.Release, appRelease.Repo)
	return appRelease.Release, err
}

func (a *AppInterface) ReloadAppReleaseResource(ctx context.Context, appReleaseResource *biz.AppReleaseResource) (*common.Msg, error) {
	err := a.appUC.ReloadAppReleaseResource(ctx, appReleaseResource)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (a *AppInterface) DeleteAppRelease(ctx context.Context, appRelease *biz.AppRelease) (*common.Msg, error) {
	err := a.appUC.DeleteAppRelease(ctx, appRelease)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (a *AppInterface) AddAppRepo(ctx context.Context, appRepo *biz.AppRepo) (*biz.AppRepo, error) {
	err := a.appUC.AddAppRepo(ctx, appRepo)
	if err != nil {
		return nil, err
	}
	return appRepo, nil
}

func (a *AppInterface) GetAppsByRepo(ctx context.Context, appRepo *biz.AppRepo) (*appApi.AppItems, error) {
	apps, err := a.appUC.GetAppsByRepo(ctx, appRepo)
	if err != nil {
		return nil, err
	}
	appItems := &appApi.AppItems{Apps: apps}
	return appItems, nil
}

func (a *AppInterface) GetAppDetailByRepo(ctx context.Context, appDetailByRepoReq *appApi.GetAppDetailByRepoReq) (*biz.App, error) {
	app, err := a.appUC.GetAppDetailByRepo(ctx, appDetailByRepoReq.Repo, appDetailByRepoReq.AppName, appDetailByRepoReq.Version)
	if err != nil {
		return nil, err
	}
	return app, nil
}
