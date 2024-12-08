use std::io::{self, Write};

#[derive(Debug, Default)]
pub struct LimitedBuffer {
    buffer: Vec<u8>,
    max_size: usize,
    try_write_size: usize,
}

impl LimitedBuffer {
    pub fn new(max_size: usize) -> Self {
        LimitedBuffer {
            buffer: Vec::with_capacity(max_size),
            max_size,
            try_write_size: 0,
        }
    }

    pub fn get_data(&self) -> &[u8] {
        &self.buffer
    }

    pub fn get_try_write_size(&self) -> usize {
        self.try_write_size
    }

    /// Consumes the LimitedBuffer and moves the buffer out
    pub fn into_buffer(self) -> Vec<u8> {
        self.buffer
    }
}

impl Write for LimitedBuffer {
    fn write(&mut self, buf: &[u8]) -> io::Result<usize> {
        self.try_write_size = self.buffer.len() + buf.len();
        if self.try_write_size > self.max_size {
            return Err(io::Error::new(
                io::ErrorKind::Other,
                "Write exceeds the buffer limit",
            ));
        }
        self.buffer.extend_from_slice(buf);
        Ok(buf.len())
    }

    fn flush(&mut self) -> io::Result<()> {
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_write_within_limit() {
        let mut buffer = LimitedBuffer::new(10);
        assert_eq!(buffer.write(b"hello").unwrap(), 5);
        assert_eq!(buffer.get_data(), b"hello");
    }

    #[test]
    fn test_write_exceeds_limit() {
        let mut buffer = LimitedBuffer::new(5);
        assert!(buffer.write(b"hello").is_ok());
        assert!(buffer.write(b"world").is_err());
        assert_eq!(buffer.get_try_write_size(), 10usize);
    }

    #[test]
    fn test_into_buffer() {
        let mut buffer = LimitedBuffer::new(10);
        buffer.write(b"hello").unwrap();
        let data = buffer.into_buffer();
        assert_eq!(data, b"hello");
    }
}
