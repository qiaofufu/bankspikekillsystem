package models

const (
	AuditPassed = 2  // 通过审核
	WaitAudit   = 1  // 等待审核
	AuditFailed = -1 // 审核失败
)

type Audit struct {
	AuditStatus  int    `json:"audit_status"`  // 审核状态
	AuditMessage string `json:"audit_message"` // 审核消息
}
