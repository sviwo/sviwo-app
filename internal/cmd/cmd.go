package cmd

import (
	"context"
	"sviwo/internal/consts"
	"sviwo/internal/controller"
	"sviwo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  consts.ProjectName,
		Usage: consts.ProjectUsage,
		Brief: consts.ProjectBrief,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// 启动管理后台gtoken
			gfAdminToken, err := StartBackendGToken()
			if err != nil {
				return err
			}
			//管理后台路由组
			s.Group("/backend", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().CORS,
					service.Middleware().Ctx,
					service.Middleware().ResponseHandler,
				)
				//不需要登录的路由组绑定
				group.Bind(
					controller.Admin.Create, // 管理员
					controller.Login,        // 登录
				)
				//需要登录的路由组绑定
				group.Group("/", func(group *ghttp.RouterGroup) {
					err := gfAdminToken.Middleware(ctx, group)
					if err != nil {
						panic(err)
					}
					group.Bind(
						controller.Role,         // 角色
						controller.Permission,   // 权限
						controller.Admin.List,   // 管理员
						controller.Admin.Update, // 管理员
						controller.Admin.Delete, // 管理员
						controller.Admin.Info,   // 查询当前管理员信息
					)
				})
			})
			//---------------------华丽的分割线-------------------
			// 启动前台项目gtoken
			frontendToken, err := StartFrontendGToken()
			if err != nil {
				return err
			}
			//前台项目路由组
			s.Group("/frontend", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().CORS,
					service.Middleware().Ctx,
					service.Middleware().ResponseHandler,
				)
				//不需要登录的路由组绑定
				group.Bind(
					controller.User.Register, //用户注册
				)
				//需要登录鉴权的路由组
				group.Group("/", func(group *ghttp.RouterGroup) {
					err := frontendToken.Middleware(ctx, group)
					if err != nil {
						return
					}
					//需要登录鉴权的接口放到这里
					group.Bind(
						controller.User.Info,           //当前登录用户的信息
						controller.User.UpdatePassword, //当前用户修改密码
					)
				})
			})
			s.SetPort(8000) //设置端口
			s.Run()
			return nil
		},
	}
)
