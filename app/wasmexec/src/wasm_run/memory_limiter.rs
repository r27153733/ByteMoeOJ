use anyhow::{Error, Result};
use wasmtime::ResourceLimiter;

pub struct MemoryLimiter {
    try_alloc_size:usize,
    max_memory: usize,
}

impl MemoryLimiter {
    pub fn new(max_memory: usize) -> Self {
        MemoryLimiter { max_memory, try_alloc_size: 0 }
    }
    pub fn get_try_alloc_size(&self) -> usize {
        self.try_alloc_size
    }
}

impl ResourceLimiter for MemoryLimiter {
    fn memory_growing(
        &mut self,
        current: usize,
        desired: usize,
        _maximum: Option<usize>,
    ) -> Result<bool> {
        self.try_alloc_size = desired;
        if desired > self.max_memory {
            Ok(false)
        } else {
            Ok(true)
        }
    }

    fn memory_grow_failed(&mut self, error: Error) -> Result<()> {
        Err(error.context("forcing a memory growth failure to be a trap"))
    }

    fn table_growing(&mut self, _current: usize, _desired: usize, _maximum: Option<usize>) -> Result<bool> {
        Ok(true)
    }

    fn table_grow_failed(&mut self, error: Error) -> Result<()> {
        Err(error.context("forcing a table growth failure to be a trap"))
    }

    fn instances(&self) -> usize {
        10000 // Limit to 10,000 instances
    }

    fn tables(&self) -> usize {
        10000 // Limit to 10,000 tables
    }

    fn memories(&self) -> usize {
        10000 // Limit to 10,000 memories
    }
}

