package models

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Definition mcp.Tool
	Handler    func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// AbuseTicket represents the AbuseTicket schema from the OpenAPI specification
type AbuseTicket struct {
	Domainip string `json:"domainIp"` // The domain or IP the suspected abuse was reported against
	Reporter string `json:"reporter"` // The shopper id of the person who reported the suspected abuse
	Ticketid string `json:"ticketId"` // Abuse ticket ID
	Closed bool `json:"closed"` // Is this abuse ticket closed?
	Closedat string `json:"closedAt"` // The date the abuse ticket was closed
	Createdat string `json:"createdAt"` // The date the abuse ticket was created
	Target string `json:"target"` // The company the suspected abuse is targeting
	TypeField string `json:"type"` // The type of abuse being reported
	Source string `json:"source"` // The single URL or IP the suspected abuse was reported against
}

// AbuseTicketCreate represents the AbuseTicketCreate schema from the OpenAPI specification
type AbuseTicketCreate struct {
	Target string `json:"target,omitempty"` // The brand/company the abuse is targeting. ie: brand name/bank name
	TypeField string `json:"type,omitempty"` // The type of abuse being reported.
	Info string `json:"info,omitempty"` // Additional information that may assist the abuse investigator. ie: server logs or email headers/body for SPAM
	Infourl string `json:"infoUrl,omitempty"` // Reporter URL if housing additional information that may assist the abuse investigator
	Intentional bool `json:"intentional,omitempty"` // Do you believe this is intentional abuse by the domain holder?
	Proxy string `json:"proxy,omitempty"` // The Proxy information required to view the abuse being reported. ie: Specific IP used, or country of IP viewing from
	Source string `json:"source,omitempty"` // The URL or IP where live abuse content is located at. ie: https://www.example.com/bad_stuff/bad.php
}

// AbuseTicketId represents the AbuseTicketId schema from the OpenAPI specification
type AbuseTicketId struct {
	U_number string `json:"u_number,omitempty"` // Abuse ticket ID
}

// AbuseTicketList represents the AbuseTicketList schema from the OpenAPI specification
type AbuseTicketList struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Ticketids []string `json:"ticketIds"` // A list of abuse ticket ids originated by this reporter.
}

// Error represents the Error schema from the OpenAPI specification
type Error struct {
	Code string `json:"code"` // Short identifier for the error, suitable for indicating the specific error within client code
	Fields []ErrorField `json:"fields,omitempty"` // List of the specific fields, and the errors found with their contents
	Message string `json:"message,omitempty"` // Human-readable, English description of the error
	Stack []string `json:"stack,omitempty"` // Stack trace indicating where the error occurred.<br/> NOTE: This attribute <strong>MAY</strong> be included for Development and Test environments. However, it <strong>MUST NOT</strong> be exposed from OTE nor Production systems.
}

// ErrorField represents the ErrorField schema from the OpenAPI specification
type ErrorField struct {
	Code string `json:"code"` // Short identifier for the error, suitable for indicating the specific error within client code
	Message string `json:"message,omitempty"` // Human-readable, English description of the problem with the contents of the field
	Path string `json:"path"` // 1) JSONPath referring to the field within the data containing an error<br/>or<br/>2) JSONPath referring to an object containing an error
	Pathrelated string `json:"pathRelated,omitempty"` // JSONPath referring to the field on the object referenced by `path` containing an error
}

// Pagination represents the Pagination schema from the OpenAPI specification
type Pagination struct {
	Total int `json:"total,omitempty"` // Number of records available
	First string `json:"first,omitempty"` // Optional link to first list of results
	Last string `json:"last,omitempty"` // Optional link to last list of results
	Next string `json:"next,omitempty"` // Optional link to next list of results
	Previous string `json:"previous,omitempty"` // Optional link to previous list of results
}
