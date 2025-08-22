package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/mcp-server/mcp-server/config"
	"github.com/mcp-server/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func CreateticketHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		// Create properly typed request body using the generated schema
		var requestBody models.AbuseTicketCreate
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/v1/abuse/tickets", cfg.BaseURL)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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
		var result models.AbuseTicketId
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

func CreateCreateticketTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_v1_abuse_tickets",
		mcp.WithDescription("Create a new abuse ticket"),
		mcp.WithString("source", mcp.Description("Input parameter: The URL or IP where live abuse content is located at. ie: https://www.example.com/bad_stuff/bad.php")),
		mcp.WithString("target", mcp.Description("Input parameter: The brand/company the abuse is targeting. ie: brand name/bank name")),
		mcp.WithString("type", mcp.Description("Input parameter: The type of abuse being reported.")),
		mcp.WithString("info", mcp.Description("Input parameter: Additional information that may assist the abuse investigator. ie: server logs or email headers/body for SPAM")),
		mcp.WithString("infoUrl", mcp.Description("Input parameter: Reporter URL if housing additional information that may assist the abuse investigator")),
		mcp.WithBoolean("intentional", mcp.Description("Input parameter: Do you believe this is intentional abuse by the domain holder?")),
		mcp.WithString("proxy", mcp.Description("Input parameter: The Proxy information required to view the abuse being reported. ie: Specific IP used, or country of IP viewing from")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    CreateticketHandler(cfg),
	}
}
