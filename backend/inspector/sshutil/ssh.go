package sshutil

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client SSH 客户端封装
type Client struct {
	host string
	port int
	user string
	conn *ssh.Client
}

// RunResult 命令执行结果
type RunResult struct {
	Stdout string
	Stderr string
	ExitCode int
}

// NewClient 创建 SSH 客户端并连接
func NewClient(host string, port int, user, password, keyFile string) (*Client, error) {
	var authMethods []ssh.AuthMethod

	// 优先使用密钥认证
	if keyFile != "" {
		key, err := os.ReadFile(keyFile)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
		// 也尝试带密码的密钥
		if len(authMethods) == 0 {
			// 如果密钥文件读取失败，回退到密码认证
			if password != "" {
				authMethods = append(authMethods, ssh.Password(password))
			}
		}
	}

	// 密码认证（如果没有密钥或密钥失败）
	if password != "" && len(authMethods) == 0 {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// 如果没有任何认证方式，尝试默认密钥
	if len(authMethods) == 0 {
		home, _ := os.UserHomeDir()
		defaultKey := home + "/.ssh/id_rsa"
		if key, err := os.ReadFile(defaultKey); err == nil {
			if signer, err := ssh.ParsePrivateKey(key); err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
		// 再尝试 id_ed25519
		if len(authMethods) == 0 {
			defaultKey = home + "/.ssh/id_ed25519"
			if key, err := os.ReadFile(defaultKey); err == nil {
				if signer, err := ssh.ParsePrivateKey(key); err == nil {
					authMethods = append(authMethods, ssh.PublicKeys(signer))
				}
			}
		}
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("无法获取 SSH 认证凭据（请提供密码或密钥文件）")
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境应使用 KnownHosts
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH 连接 %s 失败: %w", addr, err)
	}

	return &Client{
		host: host,
		port: port,
		user: user,
		conn: conn,
	}, nil
}

// Run 执行命令并返回合并输出
func (c *Client) Run(cmd string) (string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建 session 失败: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		// 命令执行失败也返回输出
		return string(output), err
	}
	return string(output), nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Addr 返回地址标识
func (c *Client) Addr() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}
