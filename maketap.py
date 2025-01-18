#!/usr/bin/env python3
import struct
import sys
from pathlib import Path

def calculate_checksum(data):
    """Calculate ZX Spectrum checksum (XOR of all bytes)"""
    return reduce(lambda x, y: x ^ y, data)

def create_header_block(filename, length, start_address=32768):
    """Create a TAP header block for CODE (bytes) format
    
    Args:
        filename: Program name (max 10 chars)
        length: Length of the data
        start_address: Memory address where code should be loaded (default: 32768/0x8000)
    """
    # Ensure filename is exactly 10 bytes, pad with spaces
    filename_bytes = filename.encode('ascii')[:10].ljust(10, b' ')
    
    # Header block is always 0x13 (19) bytes long
    header = struct.pack('<H', 0x13)  # Length of header block
    
    # Construct header data
    header_data = bytearray([
        0x00,        # Flag byte (0x00 for header)
        0x03,        # Type (0x03 for bytes/code)
    ])
    header_data.extend(filename_bytes)  # Filename (10 bytes)
    header_data.extend(struct.pack('<H', length))  # Length of data
    header_data.extend(struct.pack('<H', start_address))  # Start address
    header_data.extend([0x00, 0x00])  # Reserved bytes
    
    # Calculate and append checksum
    checksum = calculate_checksum(header_data)
    header_data.append(checksum)
    
    return header + header_data

def create_data_block(data):
    """Create a TAP data block containing the actual code/bytes
    
    Args:
        data: Binary data to store
    """
    # Data block length is data length + 2 (flag and checksum)
    length = len(data) + 2
    block = struct.pack('<H', length)  # Block length
    
    # Construct data block
    data_block = bytearray([0xFF])  # Flag byte (0xFF for data)
    data_block.extend(data)  # Actual data
    
    # Calculate and append checksum
    checksum = calculate_checksum(data_block)
    data_block.append(checksum)
    
    return block + data_block

def binary_to_tap(input_file, output_file, name=None, start_address=32768):
    """Convert a binary file to ZX Spectrum TAP format
    
    Args:
        input_file: Path to input binary file
        output_file: Path to output TAP file
        name: Name for the code block (max 10 chars)
        start_address: Memory address where code should be loaded
    """
    # Read input file
    with open(input_file, 'rb') as f:
        data = f.read()
    
    # Use input filename if no name provided
    if name is None:
        name = Path(input_file).stem[:10]
    
    # Create header and data blocks
    header_block = create_header_block(name, len(data), start_address)
    data_block = create_data_block(data)
    
    # Write TAP file
    with open(output_file, 'wb') as f:
        f.write(header_block)
        f.write(data_block)

if __name__ == '__main__':
    import argparse
    from functools import reduce
    
    parser = argparse.ArgumentParser(description='Convert binary file to ZX Spectrum TAP format')
    parser.add_argument('input', help='Input binary file')
    parser.add_argument('output', help='Output TAP file')
    parser.add_argument('--name', help='Name for code block (max 10 chars)')
    parser.add_argument('--address', type=int, default=32768,
                      help='Start address (default: 32768)')
    
    args = parser.parse_args()
    
    try:
        binary_to_tap(args.input, args.output, args.name, args.address)
        print(f"Successfully converted {args.input} to {args.output}")
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)