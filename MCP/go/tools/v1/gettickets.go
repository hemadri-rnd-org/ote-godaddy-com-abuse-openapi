package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mcp-server/mcp-server/config"
	"github.com/mcp-server/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetticketsHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["type"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("type=%v", val))
		}
		if val, ok := args["closed"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("closed=%v", val))
		}
		if val, ok := args["sourceDomainOrIp"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("sourceDomainOrIp=%v", val))
		}
		if val, ok := args["target"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("target=%v", val))
		}
		if val, ok := args["createdStart"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("createdStart=%v", val))
		}
		if val, ok := args["createdEnd"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("createdEnd=%v", val))
		}
		if val, ok := args["limit"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("limit=%v", val))
		}
		if val, ok := args["offset"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("offset=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/v1/abuse/tickets%s", cfg.BaseURL, queryString)
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
		var result models.AbuseTicketList
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

func CreateGetticketsTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_v1_abuse_tickets",
		mcp.WithDescription("List all abuse tickets ids that match user provided filters"),
		mcp.WithString("type", mcp.Description("The type of abuse.")),
		mcp.WithBoolean("closed", mcp.Description("Is this abuse ticket closed?")),
		mcp.WithString("sourceDomainOrIp", mcp.Description("The domain name or ip address the abuse originated from")),
		mcp.WithString("target", mcp.Description("The brand/company the abuse is targeting. ie: brand name/bank name")),
		mcp.WithString("createdStart", mcp.Description("The earliest abuse ticket creation date to pull abuse tickets for")),
		mcp.WithString("createdEnd", mcp.Description("The latest abuse ticket creation date to pull abuse tickets for")),
		mcp.WithNumber("limit", mcp.Description("Number of abuse ticket numbers to return.")),
		mcp.WithNumber("offset", mcp.Description("The earliest result set record number to pull abuse tickets for")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetticketsHandler(cfg),
	}
}
