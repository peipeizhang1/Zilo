// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package workflows

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zilo/backend/api/internal/logic/workflows"
	"zilo/backend/api/internal/svc"
	"zilo/backend/api/internal/types"
)

func CreateWorkflowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateWorkflowReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := workflows.NewCreateWorkflowLogic(r.Context(), svcCtx)
		resp, err := l.CreateWorkflow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
