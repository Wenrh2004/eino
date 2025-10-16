/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package schema

import (
	"bytes"
	"encoding/gob"

	"github.com/eino-contrib/jsonschema"
)

type ContentBlockType string

const (
	ContentBlockTypeMessage                 ContentBlockType = "message"
	ContentBlockTypeReasoning               ContentBlockType = "reasoning"
	ContentBlockTypeToolCall                ContentBlockType = "tool_call"
	ContentBlockTypeToolCallOutput          ContentBlockType = "tool_call_output"
	ContentBlockTypeMCPListTools            ContentBlockType = "mcp_list_tools"
	ContentBlockTypeMCPToolApprovalRequest  ContentBlockType = "mcp_tool_approval_request"
	ContentBlockTypeMCPToolApprovalResponse ContentBlockType = "mcp_tool_approval_response"
)

type AgenticResponse struct {
	ID           string
	FinishReason *FinishReason
	Usage        *TokenUsageMeta

	Blocks []*ContentBlock
}

// Serialize 由于直接使用 json marshal ，然后再 unmarshal 会丢失 extra 中的类型信息。
// 需要使用 gob 序列化。
func (r *AgenticResponse) Serialize() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(r)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (r *AgenticResponse) Deserialize(data []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(r)
}

type FinishStatus string

const (
	FinishStatusCompleted  FinishStatus = "completed"
	FinishStatusIncomplete FinishStatus = "incomplete"
)

type FinishReason struct {
	Status FinishStatus
	Reason string
}

type TokenUsageMeta struct {
	InputTokens         int64
	InputTokensDetails  InputTokensUsageDetails
	OutputTokens        int64
	OutputTokensDetails OutputTokensUsageDetails
	TotalTokens         int64
}

type InputTokensUsageDetails struct {
	CachedTokens int64
}

type OutputTokensUsageDetails struct {
	ReasoningTokens int64
}

type ContentBlock struct {
	Type ContentBlockType

	Message                 *ContentBlockMessage
	Reasoning               *ContentBlockReasoning
	ToolCall                *ContentBlockToolCall
	ToolCallOutput          *ContentBlockToolCallOutput
	MCPListTools            *ContentBlockMCPListTools
	MCPToolApprovalRequest  *ContentBlockMCPToolApprovalRequest
	MCPToolApprovalResponse *ContentBlockMCPToolApprovalResponse
}

// ContentBlockMessage system、user 、assistant message ，不包含 tool message
type ContentBlockMessage struct {
	Index *int

	Role RoleType

	InputText                string
	UserInputMultiContent    []*AgenticMessageInputPart
	AssistantGenMultiContent []*AgenticMessageOutputPart

	// 一些模型特化的参数需要存储在 Extra
	// 例如 openai output message 中的 status 和 id 等
	// 待解决问题
	// 1. 用户序列化 message 后再反序列化类型丢失问题
	Extra map[string]any
}

type AgenticMessagePartType string

const (
	AgenticMessagePartTypeText  AgenticMessagePartType = "text"
	AgenticMessagePartTypeImage AgenticMessagePartType = "image"
	AgenticMessagePartTypeAudio AgenticMessagePartType = "audio"
	AgenticMessagePartTypeVideo AgenticMessagePartType = "video"
	AgenticMessagePartTypeFile  AgenticMessagePartType = "file"
)

type AgenticMessageInputPart struct {
	Type AgenticMessagePartType

	Text  *AgenticMessageInputText
	Image *AgenticMessageInputImage
	Audio *AgenticMessageInputAudio
	Video *AgenticMessageInputVideo
	File  *AgenticMessageInputFile
}

type AgenticMessageInputText struct {
	Content string
}

type AgenticMessageInputImage struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Detail is the quality of the image url.
	Detail ImageURLDetail

	// Extra is used to store extra information.
	Extra map[string]any
}

type AgenticMessageInputAudio struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Extra is used to store extra information.
	Extra map[string]any
}

type AgenticMessageInputVideo struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Extra is used to store extra information.
	Extra map[string]any
}

type AgenticMessageInputFile struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Name is the name of the file, used when passing the file to the model as a string.
	Name *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Extra is used to store extra information.
	Extra map[string]any
}

type AgenticMessageOutputPart struct {
	Type AgenticMessagePartType

	Text  *AgenticMessageOutputText
	Image *AgenticMessageOutputImage
	Audio *AgenticMessageOutputAudio
	Video *AgenticMessageOutputVideo
}

type AgenticMessageOutputText struct {
	Content string

	Extra map[string]any
}

type AgenticMessageOutputImage struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Extra is used to store extra information.
	Extra map[string]any
}

type AgenticMessageOutputAudio struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Extra is used to store extra information.
	Extra map[string]any
}

type AgenticMessageOutputVideo struct {
	// URL can either be a traditional URL or a special URL conforming to RFC-2397 (https://www.rfc-editor.org/rfc/rfc2397).
	// double check with model implementations for detailed instructions on how to use this.
	URL *string

	// Base64Data represents the binary data in Base64 encoded string format.
	Base64Data *string

	// MIMEType is the mime type , eg."image/png",""audio/wav" etc.
	MIMEType string

	// Extra is used to store extra information.
	Extra map[string]any
}

type ContentBlockReasoning struct {
	Index        *int
	SummaryIndex *int

	Summary          []*ReasoningSummary
	EncryptedContent string

	Extra map[string]any
}

type ReasoningSummary struct {
	Text  string
	Extra map[string]any
}

type ToolCallType string

const (
	ToolCallTypeCustom ToolCallType = "custom_tool_call"
	ToolCallTypeMCP    ToolCallType = "mcp_tool_call"
)

type ContentBlockToolCall struct {
	Index *int

	Type      ToolCallType
	ID        string
	Name      string
	Arguments string

	Extra map[string]any
}

type ToolCallOutputType string

const (
	ToolCallOutputTypeCustom ToolCallOutputType = "custom_tool_call_output"
	ToolCallOutputTypeMCP    ToolCallOutputType = "mcp_tool_call_output"
)

type ContentBlockToolCallOutput struct {
	Index *int

	Type       ToolCallOutputType
	ToolCallID string
	ToolName   string

	CustomTool *ToolCallOutputCustom
	MCPTool    *ToolCallOutputMCP
}

type ToolCallOutputCustom struct {
	Content string
}

type MCPToolCallStatus string

const (
	MCPToolCallStatusSuccess MCPToolCallStatus = "success"
	MCPToolCallStatusError   MCPToolCallStatus = "error"
)

type ToolCallOutputMCP struct {
	Content string

	ApprovalRequestID string

	Status MCPToolCallStatus
	Error  string

	Extra map[string]any
}

type ContentBlockMCPListTools struct {
	// The label of the MCP server.
	ServerLabel string
	// The tools available on the server.
	Tools []MCPListToolsItem
	// Error message if the server could not list tools.
	Error string
}

type MCPListToolsItem struct {
	// The name of the tool.
	Name string
	// The description of the tool.
	Description string
	// The JSON schema describing the tool's input.
	InputSchema *jsonschema.Schema
}

type ContentBlockMCPToolApprovalRequest struct {
	// The name of the tool to run.
	Name string
	// A JSON string of arguments for the tool.
	Arguments string
	// The label of the MCP server making the request.
	ServerLabel string
}

type ContentBlockMCPToolApprovalResponse struct {
	// The ID of the approval request being answered.
	ApprovalRequestID string
	// Whether the request was approved.
	Approve bool
	// Optional reason for the decision.
	Reason string
}
