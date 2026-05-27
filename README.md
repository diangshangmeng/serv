# 1核4G服务器精准资源分配与全栈优化方案
针对您的1vCPU/4GB内存/50GB硬盘配置，我将提供**严格按1G:1G:1G:1G**资源划分的生产级优化方案，确保MySQL、Redis、Go API和系统稳定运行且互不干扰。

## 一、系统全局优化（预留1GB内存）
首先进行系统层面的基础优化，为所有服务提供稳定运行环境。

### 1. 关闭不必要的服务
```bash
# 禁用Ubuntu默认的Snap服务（节省约100-200MB内存）
sudo systemctl disable --now snapd snapd.socket
sudo systemctl mask snapd

# 禁用其他非必要服务
sudo systemctl disable --now udisks2 ModemManager bluetooth avahi-daemon cups
```

### 2. Swap配置（关键！内存受限环境必备）
```bash
# 创建1GB Swap文件（与系统预留内存相等）
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# 永久生效
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# 调整Swappiness（优先使用物理内存）
echo 'vm.swappiness=10' | sudo tee -a /etc/sysctl.conf
echo 'vm.vfs_cache_pressure=50' | sudo tee -a /etc/sysctl.conf

# 立即生效
sudo sysctl -p
```

### 3. 内核与文件描述符优化
```bash
# 编辑系统限制配置
sudo nano /etc/security/limits.conf
```
添加以下内容：
```conf
* soft nofile 65535
* hard nofile 65535
mysql soft nofile 65535
mysql hard nofile 65535
redis soft nofile 65535
redis hard nofile 65535
```

```bash
# 编辑sysctl配置
sudo nano /etc/sysctl.conf
```
添加以下内容：
```conf
# 网络优化
net.core.somaxconn=1024
net.ipv4.tcp_syncookies=1
net.ipv4.tcp_fin_timeout=30
net.ipv4.tcp_keepalive_time=120
net.ipv4.tcp_keepalive_intvl=30
net.ipv4.tcp_keepalive_probes=3

# 内存优化
vm.overcommit_memory=1
vm.dirty_ratio=10
vm.dirty_background_ratio=5
```

```bash
# 立即生效
sudo sysctl -p
```

## 二、MySQL 8.0 1GB内存精准优化
**核心原则**：将InnoDB缓冲池控制在640MB左右，总内存占用严格限制在1GB以内。

### 1. 完整优化配置文件
```bash
# 备份原配置
sudo cp /etc/mysql/mysql.conf.d/mysqld.cnf /etc/mysql/mysql.conf.d/mysqld.cnf.bak

# 编辑配置文件
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

替换为以下内容：
```ini
[mysqld]
# 基础配置
pid-file        = /var/run/mysqld/mysqld.pid
socket          = /var/run/mysqld/mysqld.sock
datadir         = /var/lib/mysql
log_error       = /var/log/mysql/error.log

# 内存限制（总内存≈1GB）
innodb_buffer_pool_size = 640M          # 核心参数，占分配内存的64%
innodb_buffer_pool_instances = 1        # 1核CPU只需1个实例
innodb_log_buffer_size = 16M
key_buffer_size = 8M
sort_buffer_size = 256K
read_buffer_size = 128K
read_rnd_buffer_size = 256K
join_buffer_size = 256K
thread_stack = 192K
tmp_table_size = 32M
max_heap_table_size = 32M

# 连接控制
max_connections = 100                   # 1核CPU足够支撑
wait_timeout = 600
interactive_timeout = 600

# InnoDB优化
innodb_flush_log_at_trx_commit = 2      # 平衡性能与安全
innodb_flush_method = O_DIRECT
innodb_file_per_table = 1
innodb_log_file_size = 64M              # 日志文件大小
innodb_log_files_in_group = 2
innodb_autoinc_lock_mode = 2
innodb_read_io_threads = 2
innodb_write_io_threads = 2
innodb_io_capacity = 200                # 适配普通SSD/HDD

# 性能与安全
performance_schema = OFF                # 关闭性能模式（节省大量内存）
skip_name_resolve = 1                   # 禁用DNS解析
default_authentication_plugin = mysql_native_password

# 日志配置
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2
log_queries_not_using_indexes = 1
```

### 2. 应用配置并验证
```bash
# 重启MySQL
sudo systemctl restart mysql

# 验证内存占用（应在800-950MB之间）
ps aux | grep mysqld
```

### 3. 重要安全加固
```sql
-- 登录MySQL执行
sudo mysql -u root -p

-- 删除匿名用户
DROP USER IF EXISTS ''@'localhost';
DROP USER IF EXISTS ''@'%';

-- 禁止root远程登录
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');

-- 刷新权限
FLUSH PRIVILEGES;
```

## 三、Redis 7.0 1GB内存精准优化
**核心原则**：设置硬内存上限1GB，启用合理的淘汰策略，优化持久化性能。

### 1. 完整优化配置文件
```bash
# 备份原配置
sudo cp /etc/redis/redis.conf /etc/redis/redis.conf.bak

# 编辑配置文件
sudo nano /etc/redis/redis.conf
```

修改以下关键参数：
```ini
# 内存限制（硬上限）
maxmemory 512mb
maxmemory-policy allkeys-lru           # 优先淘汰最近最少使用的键
maxmemory-samples 5

# 持久化优化（平衡性能与数据安全）
save 900 1                             # 900秒内至少1个键发生变化则保存
save 300 10                            # 300秒内至少10个键发生变化则保存
save 60 10000                          # 60秒内至少10000个键发生变化则保存
stop-writes-on-bgsave-error no
rdbcompression yes
rdbchecksum yes

# AOF配置（可选，根据数据重要性选择）
appendonly no                          # 如需要更高数据安全性改为yes
appendfsync everysec                   # 每秒同步一次（推荐）
no-appendfsync-on-rewrite yes
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# 网络优化
bind 127.0.0.1 ::1                     # 仅本地监听（如需远程改为0.0.0.0）
port 6379
tcp-backlog 511
timeout 300
tcp-keepalive 300

# 性能优化
databases 16                           # 减少数据库数量
save ""                                # 如完全不需要持久化可取消注释
lazyfree-lazy-eviction yes
lazyfree-lazy-expire yes
lazyfree-lazy-server-del yes
replica-lazy-flush yes

# 安全配置
requirepass your_strong_redis_password # 必须设置强密码
rename-command FLUSHDB ""              # 禁用危险命令
rename-command FLUSHALL ""
rename-command KEYS ""
rename-command DEBUG ""
rename-command CONFIG ""
```

### 2. 应用配置并验证
```bash
# 重启Redis
sudo systemctl restart redis-server

# 验证内存限制
redis-cli
AUTH your_strong_redis_password
INFO memory
```
查看`used_memory_human`字段，确保不会超过1GB。

## 四、Go API服务 1GB内存优化
**核心原则**：利用Go 1.19+的内存限制特性，精确控制GC行为，避免内存溢出。

### 1. 编译优化
在编译Go程序时添加以下参数，减小二进制文件大小并提高性能：
```bash
go build -ldflags="-s -w" -o api-server main.go
```

### 2. 运行时环境变量优化
创建启动脚本`start-api.sh`：
```bash
#!/bin/bash

# 内存限制（软限制1GB，硬限制1.1GB）
export GOMEMLIMIT=1GiB
export GOGC=50                          # 降低GC阈值，更频繁回收内存
export GOMAXPROCS=1                     # 1核CPU设置为1

# 启动服务
./api-server
```

```bash
# 添加执行权限
chmod +x start-api.sh
```

### 3. Systemd服务配置（推荐）
创建systemd服务文件`/etc/systemd/system/api-server.service`：
```ini
[Unit]
Description=Go API Server
After=network.target mysql.service redis-server.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/api-server
ExecStart=/opt/api-server/start-api.sh

# 资源限制（强制1GB内存上限）
MemoryMax=1G
MemoryHigh=900M
CPUQuota=100%

# 自动重启
Restart=always
RestartSec=5

# 日志配置
StandardOutput=journal+console
StandardError=journal+console
SyslogIdentifier=api-server

[Install]
WantedBy=multi-user.target
```

```bash
# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable --now api-server
```

### 4. Go代码层面优化建议
- 使用`sync.Pool`复用对象，减少GC压力
- 避免使用全局变量和大对象
- 及时关闭不需要的文件和网络连接
- 设置HTTP服务器超时：
  ```go
  server := &http.Server{
      Addr:         ":8080",
      ReadTimeout:  5 * time.Second,
      WriteTimeout: 10 * time.Second,
      IdleTimeout:  120 * time.Second,
      MaxHeaderBytes: 1 << 20, // 1MB
  }
  ```

## 五、服务资源隔离与监控
### 1. 验证所有服务内存占用
```bash
# 查看各服务内存使用情况
ps aux --sort=-%mem | head -10

# 更详细的内存统计
free -h
```
预期结果：
- MySQL：约800-950MB
- Redis：约900-1000MB（含数据）
- Go API：约500-900MB（根据并发量）
- 系统：约500-800MB

### 2. 安装简单监控工具
```bash
sudo apt install htop iotop -y
```

使用`htop`可以实时查看各服务的CPU和内存使用情况。

## 六、最终资源分配总结
| 服务 | 分配内存 | 实际占用范围 | 核心优化点 |
|------|----------|--------------|------------|
| 系统 | 1GB | 500-800MB | 关闭不必要服务，优化Swap |
| MySQL | 1GB | 800-950MB | 限制InnoDB缓冲池，关闭性能模式 |
| Redis | 1GB | 900-1000MB | 设置硬内存上限，启用LRU淘汰 |
| Go API | 1GB | 500-900MB | GOMEMLIMIT限制，优化GC |

## 七、注意事项
1. **硬盘监控**：50GB硬盘需定期清理日志文件，建议配置日志轮转
2. **备份策略**：MySQL每天备份一次，Redis根据数据重要性定期备份
3. **安全更新**：定期更新系统和软件包
4. **性能监控**：如发现内存不足，可适当降低MySQL的`max_connections`或Redis的`maxmemory`


# 基于Cloudflare Zero Trust的安全架构调整方案
完全符合您的要求：**仅通过Cloudflare Zero Trust暴露API，MySQL和Redis彻底关闭外部访问，仅保留22端口SSH**。这种架构是目前最安全的自托管API部署方式，完全隐藏服务器真实IP，内置DDoS防护和WAF。

## 一、防火墙终极配置（只开放22端口）
**立即关闭所有不必要端口**，这是安全的第一道防线。

```bash
# 重置所有防火墙规则
sudo ufw reset

# 只允许SSH访问（22端口）
sudo ufw allow 22/tcp

# 拒绝所有其他入站连接
sudo ufw default deny incoming

# 允许所有出站连接
sudo ufw default allow outgoing

# 启用防火墙
sudo ufw enable

# 验证防火墙状态
sudo ufw status
```

**验证输出应该只有这一行**：
```
Status: active

To                         Action      From
--                         ------      ----
22/tcp                     ALLOW       Anywhere
22/tcp (v6)                ALLOW       Anywhere (v6)
```

⚠️ **重要提醒**：执行此操作前确保你当前的SSH连接不会断开，否则将无法远程登录服务器。

## 二、MySQL和Redis彻底关闭外部访问
确保这两个服务**只能通过本地回环地址127.0.0.1访问**，即使防火墙被绕过也无法被外部攻击。

### 1. MySQL本地访问加固
```bash
# 编辑MySQL配置文件
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

确保以下配置（**必须**）：
```ini
# 只监听本地回环地址
bind-address = 127.0.0.1
mysqlx-bind-address = 127.0.0.1

# 禁用DNS解析（进一步提升性能和安全）
skip_name_resolve = 1
```

```bash
# 重启MySQL应用配置
sudo systemctl restart mysql

# 验证监听地址
sudo netstat -tulpn | grep mysql
```

**验证输出**：
```
tcp        0      0 127.0.0.1:3306          0.0.0.0:*               LISTEN      1234/mysqld
tcp        0      0 127.0.0.1:33060         0.0.0.0:*               LISTEN      1234/mysqld
```

### 2. Redis本地访问加固
```bash
# 编辑Redis配置文件
sudo nano /etc/redis/redis.conf
```

确保以下配置（**必须**）：
```ini
# 只监听本地回环地址（删除所有其他bind行）
bind 127.0.0.1 ::1

# 禁用保护模式（因为只监听本地，保护模式已无必要）
protected-mode no

# 保留密码认证（防止本地提权攻击）
requirepass your_strong_redis_password
```

```bash
# 重启Redis应用配置
sudo systemctl restart redis-server

# 验证监听地址
sudo netstat -tulpn | grep redis
```

**验证输出**：
```
tcp        0      0 127.0.0.1:6379          0.0.0.0:*               LISTEN      5678/redis-server 1
tcp6       0      0 ::1:6379                :::*                    LISTEN      5678/redis-server 1
```

## 三、Cloudflare Zero Trust隧道配置
通过Cloudflare Tunnel将本地API服务安全暴露到公网，**不需要开放任何业务端口**。

### 1. 安装Cloudflare Tunnel客户端
```bash
# 添加Cloudflare GPG密钥
curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null

# 添加Cloudflare软件源
echo "deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/cloudflared.list

# 更新软件包并安装cloudflared
sudo apt update && sudo apt install cloudflared -y
```

### 2. 认证并创建隧道
```bash
# 登录Cloudflare账户
sudo cloudflared tunnel login
```

执行后会生成一个链接，在浏览器中打开并授权你的域名。

```bash
# 创建隧道（替换为你的隧道名称）
sudo cloudflared tunnel create api-tunnel

# 验证隧道创建成功
sudo cloudflared tunnel list
```

### 3. 配置隧道路由
```bash
# 将域名指向隧道（替换为你的域名和隧道ID）
sudo cloudflared tunnel route dns api-tunnel api.yourdomain.com

# 创建隧道配置文件
sudo mkdir -p /etc/cloudflared
sudo nano /etc/cloudflared/config.yml
```

写入以下配置：
```yaml
tunnel: 你的隧道ID
credentials-file: /root/.cloudflared/你的隧道ID.json

ingress:
  # 将所有请求转发到本地运行的Go API服务
  - hostname: api.yourdomain.com
    service: http://localhost:8080
  
  # 拒绝所有其他请求
  - service: http_status:404
```

### 4. 配置Systemd服务自动启动
```bash
# 安装服务
sudo cloudflared service install

# 启动并启用服务
sudo systemctl enable --now cloudflared

# 验证服务状态
sudo systemctl status cloudflared
```

### 5. 验证隧道连接
```bash
# 查看隧道连接状态
sudo cloudflared tunnel info api-tunnel
```

如果显示"Connected"，说明隧道已成功建立。现在你可以通过`https://api.yourdomain.com`访问你的API服务了。

## 四、Go API服务调整
由于API服务现在只需要监听本地地址，我们可以进一步提升安全性。

### 1. 修改API服务监听地址
在你的Go代码中，将监听地址从`0.0.0.0:8080`改为`127.0.0.1:8080`：

```go
// 原来的代码
server := &http.Server{
    Addr: ":8080",
    // ...其他配置
}

// 修改后的代码
server := &http.Server{
    Addr: "127.0.0.1:8080", // 只监听本地回环地址
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
    MaxHeaderBytes: 1 << 20, // 1MB
}
```

### 2. 更新Systemd服务配置
确保API服务在Cloudflare Tunnel之后启动：

```bash
sudo nano /etc/systemd/system/api-server.service
```

修改`After`行：
```ini
[Unit]
Description=Go API Server
After=network.target mysql.service redis-server.service cloudflared.service
```

```bash
# 重新加载配置并重启服务
sudo systemctl daemon-reload
sudo systemctl restart api-server
```

## 五、Cloudflare Zero Trust额外安全配置
### 1. 启用WAF（Web应用防火墙）
1. 登录Cloudflare仪表板
2. 进入你的域名 -> 安全 -> WAF
3. 启用"Cloudflare Managed Ruleset"
4. 启用"SQL Injection"和"XSS"防护规则

### 2. 配置访问策略（可选但强烈推荐）
如果你只想允许特定IP或用户访问API，可以配置Zero Trust访问策略：
1. 进入Cloudflare Zero Trust仪表板 -> 访问 -> 应用
2. 添加应用 -> 自托管
3. 输入应用名称和域名
4. 创建访问策略，只允许你指定的IP地址或用户访问

### 3. 启用缓存和压缩
1. 进入Cloudflare仪表板 -> 缓存 -> 配置
2. 启用"自动压缩"
3. 配置适当的缓存规则（根据你的API响应特性）

## 六、最终安全验证
### 1. 端口扫描验证
使用在线端口扫描工具（如https://www.grc.com/shieldsup）扫描你的服务器IP，**应该只有22端口显示为开放**，其他所有端口都显示为关闭或过滤。

### 2. 本地访问验证
```bash
# 测试MySQL本地访问
mysql -u root -p -h 127.0.0.1

# 测试Redis本地访问
redis-cli -h 127.0.0.1 -a your_strong_redis_password

# 测试API本地访问
curl http://127.0.0.1:8080/health
```

### 3. 远程访问验证
```bash
# 测试API通过Cloudflare访问
curl https://api.yourdomain.com/health

# 测试直接访问服务器IP的API（应该失败）
curl http://你的服务器IP:8080/health

# 测试直接访问MySQL（应该失败）
mysql -u root -p -h 你的服务器IP

# 测试直接访问Redis（应该失败）
redis-cli -h 你的服务器IP
```

## 七、架构优势总结
1. **完全隐藏真实IP**：攻击者无法直接攻击你的服务器
2. **零端口暴露**：除了22端口SSH外，没有任何业务端口开放
3. **内置DDoS防护**：Cloudflare提供全球DDoS防护能力
4. **免费SSL证书**：自动生成和续期SSL证书
5. **细粒度访问控制**：可以基于IP、用户、设备等控制API访问
6. **全球加速**：Cloudflare全球CDN网络加速API访问

## 八、注意事项
1. **SSH安全**：虽然只开放了22端口，但仍建议使用密钥认证而非密码认证
2. **隧道监控**：定期检查Cloudflare Tunnel的运行状态
3. **日志管理**：Cloudflare提供详细的访问日志，建议定期查看
4. **备份策略**：继续保持MySQL和Redis的定期备份
