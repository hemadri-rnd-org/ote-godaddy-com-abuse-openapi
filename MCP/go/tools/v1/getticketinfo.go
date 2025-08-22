package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mcp-server/mcp-server/config"
	"github.com/mcp-server/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetticketinfoHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		ticketIdVal, ok := args["ticketId"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: ticketId"), nil
		}
		ticketId, ok := ticketIdVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: ticketId"), nil
		}
		url := fmt.Sprintf("%s/v1/abuse/tickets/%s", cfg.BaseURL, ticketId)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No authentication required for this endpoint
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.AbuseTicket
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetticketinfoTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_v1_abuse_tickets_ticketId",
		mcp.WithDescription("Return the abuse ticket data for a given ticket id"),
		mcp.WithString("ticketId", mcp.Required(), mcp.Description("A unique abuse ticket identifier")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetticketinfoHandler(cfg),
	}
}
