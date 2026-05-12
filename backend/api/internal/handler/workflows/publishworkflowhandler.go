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

func PublishWorkflowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflowID, err := parsePathUint64(r, "id")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		var req types.PublishReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := workflows.NewPublishWorkflowLogic(r.Context(), svcCtx)
		resp, err := l.PublishWorkflow(workflowID, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
