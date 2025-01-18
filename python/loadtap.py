#!/usr/bin/env python3
import struct
import sys
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class TAPHeader:
    type: int
    filename: str
    data_length: int
    param1: int
    param2: int
    checksum: int

    @classmethod
    def from_bytes(cls, header_data: bytes, checksum: int) -> 'TAPHeader':
        """Parse header data (without checksum) and create a TAPHeader object"""
        if len(header_data) != 17:  # Header should be exactly 17 bytes (without checksum)
            raise ValueError(f"Invalid header data length: {len(header_data)}")
            
        type_byte = header_data[0]
        filename = header_data[1:11].decode('ascii').rstrip()
        data_length = struct.unpack('<H', header_data[11:13])[0]
        param1 = struct.unpack('<H', header_data[13:15])[0]
        param2 = struct.unpack('<H', header_data[15:17])[0]
        
        return cls(type_byte, filename, data_length, param1, param2, checksum)

@dataclass
class TAPBlock:
    length: int
    flag: int
    data: bytes
    checksum: int
    header: Optional[TAPHeader] = None

def read_tap_file(filename: str) -> List[TAPBlock]:
    """Read a TAP file and return a list of its blocks"""
    blocks = []
    
    with open(filename, 'rb') as f:
        while True:
            # Read block length (2 bytes, little-endian)
            length_bytes = f.read(2)
            if not length_bytes:
                break
            
            length = struct.unpack('<H', length_bytes)[0]
            
            # Read flag byte
            flag = f.read(1)[0]
            
            # Read block data (excluding checksum)
            data = f.read(length - 2)  # -2 for flag and checksum
            
            # Read checksum
            checksum = f.read(1)[0]
            
            # Create block object
            block = TAPBlock(length, flag, data, checksum)
            
            # If this is a header block (flag = 0x00 and length = 0x13)
            if flag == 0x00 and length == 0x13:
                try:
                    block.header = TAPHeader.from_bytes(data, checksum)
                except Exception as e:
                    print(f"Warning: Failed to parse header: {e}", file=sys.stderr)
            
            blocks.append(block)
    
    return blocks

def print_block_info(block: TAPBlock, index: int):
    """Print information about a TAP block"""
    print(f"\nBlock {index}:")
    print(f"  Length: {block.length}")
    print(f"  Flag: 0x{block.flag:02X} ({'Header' if block.flag == 0x00 else 'Data'})")
    
    if block.header:
        print("  Header Information:")
        print(f"    Type: {block.header.type}")
        print(f"    Filename: {block.header.filename}")
        print(f"    Data Length: {block.header.data_length}")
        print(f"    Param1: {block.header.param1}")
        print(f"    Param2: {block.header.param2}")
    print(f"  Checksum: 0x{block.checksum:02X}")
    print(f"  Data Length: {len(block.data)} bytes")

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='Read ZX Spectrum TAP file')
    parser.add_argument('file', help='TAP file to read')
    parser.add_argument('-d', '--dump', action='store_true', help='Dump block data as hex')
    parser.add_argument('-r', '--raw', action='store_true', help='Output raw block data')
    args = parser.parse_args()
    
    try:
        blocks = read_tap_file(args.file)
        
        if args.raw:
            # Output just the raw data blocks (skip headers)
            for block in blocks:
                if block.flag != 0x00:  # If it's not a header block
                    sys.stdout.buffer.write(block.data)
        else:
            # Print information about each block
            print(f"Found {len(blocks)} blocks in {args.file}")
            
            for i, block in enumerate(blocks):
                print_block_info(block, i)
                
                if args.dump:
                    print("  Data:")
                    hex_dump = ' '.join(f'{b:02X}' for b in block.data)
                    for i in range(0, len(hex_dump), 48):
                        print(f"    {hex_dump[i:i+48]}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()