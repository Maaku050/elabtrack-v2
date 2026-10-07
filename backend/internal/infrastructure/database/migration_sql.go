package database

import (
	"errors"
	"strings"
)

// A narrow lexical guard prevents a file from committing the runner-owned
// transaction. It understands PostgreSQL quotes, dollar bodies and nested
// comments; PostgreSQL remains the SQL parser. Session SET/RESET and prepared
// statements are deliberately excluded from this transactional file convention.
func validateTransactionalSQL(sql string) error {
	first := true
	for i := 0; i < len(sql); {
		c := sql[i]
		if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			i++
			continue
		}
		if i+1 < len(sql) && sql[i:i+2] == "--" {
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			continue
		}
		if i+1 < len(sql) && sql[i:i+2] == "/*" {
			depth := 1
			i += 2
			for i < len(sql) && depth > 0 {
				if i+1 < len(sql) && sql[i:i+2] == "/*" {
					depth++
					i += 2
				} else if i+1 < len(sql) && sql[i:i+2] == "*/" {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth != 0 {
				return errors.New("unterminated comment")
			}
			continue
		}
		if c == ';' {
			first = true
			i++
			continue
		}
		escape := false
		if (c == 'E' || c == 'e') && i+1 < len(sql) && sql[i+1] == '\'' {
			escape = true
			i++
			c = sql[i]
		}
		if c == '\'' || c == '"' {
			quote := c
			i++
			closed := false
			for i < len(sql) {
				if escape && sql[i] == '\\' {
					i += 2
					continue
				}
				if sql[i] == quote {
					if i+1 < len(sql) && sql[i+1] == quote {
						i += 2
						continue
					}
					i++
					closed = true
					break
				}
				i++
			}
			if !closed {
				return errors.New("unterminated quote")
			}
			first = false
			continue
		}
		if c == '$' {
			end := i + 1
			for end < len(sql) && (sql[end] == '_' || sql[end] >= 'a' && sql[end] <= 'z' || sql[end] >= 'A' && sql[end] <= 'Z' || sql[end] >= '0' && sql[end] <= '9') {
				end++
			}
			if end < len(sql) && sql[end] == '$' {
				tag := sql[i : end+1]
				close := strings.Index(sql[end+1:], tag)
				if close < 0 {
					return errors.New("unterminated dollar body")
				}
				i = end + 1 + close + len(tag)
				first = false
				continue
			}
		}
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' {
			start := i
			for i < len(sql) && (sql[i] >= 'a' && sql[i] <= 'z' || sql[i] >= 'A' && sql[i] <= 'Z' || sql[i] >= '0' && sql[i] <= '9' || sql[i] == '_') {
				i++
			}
			word := strings.ToUpper(sql[start:i])
			if first {
				switch word {
				case "BEGIN", "START", "COMMIT", "END", "ROLLBACK", "ABORT", "SAVEPOINT", "RELEASE", "PREPARE", "SET", "RESET":
					return errors.New("runner owns transaction/session control")
				}
			}
			first = false
			continue
		}
		first = false
		i++
	}
	return nil
}
