#!/usr/bin/env python3

from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import os
from urllib.parse import urlparse
from datetime import datetime

# Global configuration
AUTH_ENABLED = False
BEARER_TOKEN = "test_token"

class FileHandler(BaseHTTPRequestHandler):
    FILE_PATH = "data.json"
    
    def get_default_data(self):
        """Generate default data structure"""
        return {
            "created": datetime.now().isoformat(),
            "data": [
                {
                    "data": "bruh"
                },
                {
                    "data2": ["string"]
                }
            ]
        }
    
    def ensure_file_exists(self):
        """Create data.json with defaults if it doesn't exist"""
        if not os.path.exists(self.FILE_PATH):
            default_data = self.get_default_data()
            with open(self.FILE_PATH, 'w') as f:
                json.dump(default_data, f, indent=2)
            print(f"Created {self.FILE_PATH} with default structure")
    
    def do_GET(self):
        # Simple auth check
        if AUTH_ENABLED:
            auth_header = self.headers.get('Authorization', '')
            expected_token = f"Bearer {BEARER_TOKEN}"
            if auth_header != expected_token:
                self.send_error(401, "Unauthorized")
                return
        
        parsed_path = urlparse(self.path)
        
        if parsed_path.path == "/file":
            try:
                self.ensure_file_exists()
                
                with open(self.FILE_PATH, 'r') as f:
                    data = json.load(f)
                
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(data, indent=2).encode())
                    
            except json.JSONDecodeError:
                self.send_error(500, "Invalid JSON in file")
            except Exception as e:
                self.send_error(500, f"Server error: {str(e)}")
        else:
            self.send_error(404, "Not found")
    
    def do_POST(self):
        self.do_PUT()

    def do_PUT(self):
        # Simple auth check
        if AUTH_ENABLED:
            auth_header = self.headers.get('Authorization', '')
            expected_token = f"Bearer {BEARER_TOKEN}"
            if auth_header != expected_token:
                self.send_error(401, "Unauthorized")
                return
        
        parsed_path = urlparse(self.path)
        
        if parsed_path.path == "/file":
            try:
                content_length = int(self.headers.get('Content-Length', 0))
                if content_length == 0:
                    self.send_error(400, "No content provided")
                    return
                
                put_data = self.rfile.read(content_length)
                
                try:
                    json_data = json.loads(put_data.decode())
                except json.JSONDecodeError:
                    self.send_error(400, "Invalid JSON format")
                    return
                
                with open(self.FILE_PATH, 'w') as f:
                    json.dump(json_data, f, indent=2)
                
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(b'{"status": "success", "message": "File updated"}')
                
            except Exception as e:
                self.send_error(500, f"Server error: {str(e)}")
        else:
            self.send_error(404, "Not found")

def run_server(port=8000):
    server_address = ('', port)
    httpd = HTTPServer(server_address, FileHandler)
    
    print(f"Server running on http://localhost:{port}")
    print(f"Auth enabled: {AUTH_ENABLED}")
    print(f"Bearer token: {BEARER_TOKEN}")
    print()
    print("Endpoints:")
    print(f"  GET  http://localhost:{port}/file - Display JSON")
    print(f"  PUT  http://localhost:{port}/file - Replace JSON")
    
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down server...")
        httpd.shutdown()

if __name__ == "__main__":
    run_server()