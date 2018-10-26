import http.client
import json
import sys

crawler_url = sys.argv[1]
crawler_host = crawler_url[crawler_url.find('://')+3:]
crawl_seed_url = sys.argv[2]
connection = None

def url_keymap(url):
    filename = url[url.rfind("/")+1:]
    return {filename: url}

try:
    if crawler_url.startswith('https'):
        connection = http.client.HTTPSConnection(crawler_host)
    else:
        connection = http.client.HTTPConnection(crawler_host)

        headers = {'Content-type': 'application/json'}

        crawl_request = {"repository": "file",
                         "repositoryConfig": "/home/vijay/Desktop/repository",
                         "resources": url_keymap(crawl_seed_url),
                         "context": {}}
        json_req = json.dumps(crawl_request)
        
        connection.request('POST', '/crawl', json_req, headers)

        response = connection.getresponse()
        print(response.read().decode())
finally:
    if connection:
        connection.close()
