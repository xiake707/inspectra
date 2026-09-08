export const hostRules = {
    Host: [
        { 
            required: true,
            message: '请输入主机地址', 
            trigger: 'blur' 
        }, 
        { validator: 
            (_, value, callback) => { 
                const ip = /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/; 
                ip.test(value) ? callback() : callback(new Error('请输入合法的 IPv4 地址')) 
            }, 
            trigger: 'blur' 
        }
    ],
    User: [
        { 
            required: true, 
            message: '请输入 SSH 用户', 
            trigger: 'blur' 
        },
        {
            pattern: /^[A-Za-z_][A-Za-z0-9_.-]{0,19}$/,
            message: '仅允许字母、数字、_、.、-，最多 20 位',
            trigger: 'blur' 
        }
    ],
    Port: [
        { 
            required: true, 
            message: '请输入 SSH 端口', 
            trigger: 'blur' 
        }, 
        { 
            type: 'number',
            min: 1, 
            max: 65535, 
            message: '端口范围为 1 - 65535',
            trigger: 'blur' 
        }
    ],
    Password: [
        { 
            max: 20, 
            message: '密码最多 20 位', 
            trigger: 'blur' 
        }
    ],
    KeyFile: [
        {
            max: 30, 
            message: '私钥路径最多 30 位', 
            trigger: 'blur' 
        }, 
        {
            validator: (_, value, callback) => {
                if (!value || value.startsWith('/') || value.startsWith('$HOME/')) 
                    callback(); 
                else 
                    callback(new Error('需为绝对路径或 $HOME/ 开头')) 
            }, 
            trigger: 'blur' 
        }
    ]
}
