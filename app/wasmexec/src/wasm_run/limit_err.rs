use std::fmt;

pub enum LimitErr {
    FuelLimitExceeded,
    MemoryLimitExceeded,
    StdoutOutputLimitExceeded,
    StderrOutputLimitExceeded,
}

impl fmt::Display for LimitErr {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let description = match self {
            LimitErr::FuelLimitExceeded => "Fuel limit exceeded",
            LimitErr::MemoryLimitExceeded => "Memory limit exceeded",
            LimitErr::StdoutOutputLimitExceeded => "Stdout output limit exceeded",
            LimitErr::StderrOutputLimitExceeded => "Stderr output limit exceeded",
        };
        write!(f, "{}", description)
    }
}

// 为 LimitErr 实现 fmt::Debug 特征
impl fmt::Debug for LimitErr {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self) // 使用 Display 的输出
    }
}

// 为 LimitErr 实现 std::error::Error 特征
impl std::error::Error for LimitErr {}