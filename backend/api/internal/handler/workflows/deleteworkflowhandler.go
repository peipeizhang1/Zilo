// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package workflows

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zilo/backend/api/internal/logic/workflows"
	"zilo/backend/api/internal/svc"
)

func DeleteWorkflowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workflowID, err := parsePathUint64(r, "id")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := workflows.NewDeleteWorkflowLogic(r.Context(), svcCtx)
		err = l.DeleteWorkflow(workflowID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
