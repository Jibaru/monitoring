package scripts

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type ReceiveLogsReq struct {
	AppKey  string   `json:"-"`
	Logs    []string `json:"logs"`
	LogType *string  `json:"logType"`
}

type ReceiveLogsResp struct {
	Message string `json:"message"`
}

type ReceiveLogsScript struct {
	db *mongo.Database
}

func NewReceiveLogsScript(db *mongo.Database) *ReceiveLogsScript {
	return &ReceiveLogsScript{db: db}
}

func (s *ReceiveLogsScript) Exec(ctx context.Context, req ReceiveLogsReq) (*ReceiveLogsResp, error) {
	app, err := persistence.GetAppByKey(ctx, s.db, req.AppKey)
	if err != nil {
		return nil, err
	}

	logType := "json"
	if req.LogType != nil {
		logType = *req.LogType
	}

	logs := make([]persistence.Log, len(req.Logs))
	for i, rawLog := range req.Logs {
		data := s.parse(rawLog, logType)

		logs[i] = persistence.Log{
			ID:        primitive.NewObjectID(),
			AppID:     app.ID,
			Timestamp: time.Now().UTC(),
			Data:      data,
			Raw:       rawLog,
		}
	}

	err = persistence.SaveLogs(ctx, s.db, logs)
	if err != nil {
		return nil, err
	}

	return &ReceiveLogsResp{Message: "Logs recibidos"}, nil
}

func (s *ReceiveLogsScript) parse(rawLog string, logType string) map[string]any {
	switch strings.ToLower(logType) {
	case "json":
		var data map[string]any
		if err := json.Unmarshal([]byte(rawLog), &data); err != nil {
			return nil
		}
		return data
	case "xml":
		var root xmlNode
		if err := xml.Unmarshal([]byte(rawLog), &root); err != nil {
			return nil
		}
		return nodeToMap(root)
	case "apache":
		// Example:
		// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
		pattern := `^(?P<ip>\S+) \S+ \S+ \[(?P<time>[^\]]+)\] "(?P<request>[^"]+)" (?P<status>\d{3}) (?P<size>\S+)$`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(rawLog)
		if matches == nil {
			return nil
		}
		result := map[string]any{
			"ip":      matches[1],
			"time":    matches[2],
			"request": matches[3],
			"status":  matches[4],
			"size":    matches[5],
		}
		if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2]); err == nil {
			result["timestamp"] = t.UTC()
		}
		return result
	case "nginx":
		pattern := `^(?P<ip>\S+) - \S+ \[(?P<time>[^\]]+)\] "(?P<request>[^"]+)" (?P<status>\d{3}) (?P<size>\d+)$`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(rawLog)
		if matches == nil {
			return nil
		}
		result := map[string]any{
			"ip":      matches[1],
			"time":    matches[2],
			"request": matches[3],
			"status":  matches[4],
			"size":    matches[5],
		}
		if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2]); err == nil {
			result["timestamp"] = t.UTC()
		}
		return result
	case "syslog":
		// Example:
		// "Mar 30 15:04:05 hostname process: message"
		pattern := `^(?P<month>\w{3})\s+(?P<day>\d{1,2})\s+(?P<time>\d{2}:\d{2}:\d{2})\s+(?P<host>\S+)\s+(?P<process>[^:]+):\s+(?P<message>.+)$`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(rawLog)
		if matches == nil {
			return nil
		}
		result := map[string]any{
			"month":   matches[1],
			"day":     matches[2],
			"time":    matches[3],
			"host":    matches[4],
			"process": matches[5],
			"message": matches[6],
		}
		return result
	case "csv":
		fields := strings.Split(rawLog, ",")
		if len(fields) == 0 {
			return nil
		}
		data := make(map[string]any, len(fields))
		for i, field := range fields {
			key := fmt.Sprintf("field%d", i+1)
			data[key] = strings.TrimSpace(field)
		}
		return data
	case "plain":
		return map[string]any{
			"message": rawLog,
		}
	default:
		return nil
	}
}

type xmlNode struct {
	XMLName  xml.Name   `xml:""`
	Attrs    []xml.Attr `xml:"-"`
	Content  string     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
}

func nodeToMap(n xmlNode) map[string]any {
	m := map[string]any{
		"tag":     n.XMLName.Local,
		"content": strings.TrimSpace(n.Content),
	}
	if len(n.Attrs) > 0 {
		attrMap := make(map[string]string)
		for _, a := range n.Attrs {
			attrMap[a.Name.Local] = a.Value
		}
		m["attrs"] = attrMap
	}
	if len(n.Children) > 0 {
		children := make([]map[string]any, 0, len(n.Children))
		for _, child := range n.Children {
			children = append(children, nodeToMap(child))
		}
		m["children"] = children
	}
	return m
}
