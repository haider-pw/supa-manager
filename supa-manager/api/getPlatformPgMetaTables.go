package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformPgMetaTables(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Get query parameters
	includeColumns := c.Query("include_columns") == "true"
	includedSchemas := c.Query("included_schemas")

	// TODO: In Phase 3/4, connect to actual project's PostgreSQL database
	// For now, return empty array to prevent UI crashes
	// Studio expects an array of table objects with this structure:
	// {
	//   id: number,
	//   schema: string,
	//   name: string,
	//   rls_enabled: boolean,
	//   rls_forced: boolean,
	//   replica_identity: string,
	//   bytes: number,
	//   size: string,
	//   live_rows_estimate: number,
	//   dead_rows_estimate: number,
	//   comment: string | null,
	//   columns: [...] (if include_columns=true)
	// }

	a.logger.Info("Fetching pg-meta tables",
		"include_columns", includeColumns,
		"included_schemas", includedSchemas)

	// Return empty array - Studio will show "No tables found"
	c.JSON(http.StatusOK, []interface{}{})
}
