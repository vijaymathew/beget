import http.client
import json
import sys
import getopt
import fileinput

def print_usage():
    print('conroller.py <options>\n')
    print('options: \n')
    print('-h                    print this help and quit\n')
    print('-c --crawlerhost=URL  url to reach the crawler service\n')
    print('-f --urls=filename    name of file with URLs to crawl, one on each line\n')
    print('-r --repodir=PATH     full path to the directory where crawled documents are stored\n')

crawler_url = ''
crawl_seed_file = ''
crawl_repository = '.'

try:
    opts, args = getopt.getopt(argv,"hc:f:r",["crawlerhost=","urls=", "repodir="])
except getopt.GetoptError:
    print_usage()
    sys.exit(2)

for opt, arg in opts:
    if opt == '-h':
        print_usage()
        sys.exit()
    elif opt in ("-c", "--crawlerhost"):
        crawler_url = arg
    elif opt in ("-f", "--urls"):
        crawl_seed_file = arg
    elif opt in ("-r", "--repodir"):
        crawl_repository = arg

crawler_host = crawler_url[crawler_url.find('://')+3:]

def urlfile(url):
    return url[url.rfind("/")+1:]

def getconnection():
    if crawler_url.startswith('https'):
        return http.client.HTTPSConnection(crawler_host)
    else:
        return http.client.HTTPConnection(crawler_host)

headers = {'Content-type': 'application/json'}

def crawl(urls):
    resources = {}
    for url in urls:
        k = urlfile(url)
        resources[k] = url
        docs_q.put({'filename': k, 'retries': 0, 'url': url})

    request = {"repository": "file",
               "repositoryConfig": crawl_repository,
               "resources": {},
               "context": {}}
    request["resources"] = resources
    json_req = json.dumps(request)
    print("Crawl request: " + json_req)
    connection = None
    try:
        connection = getconnection()
        connection.request('POST', '/crawl', json_req, headers)
        response = connection.getresponse()
    except Exception as ex:
        print(ex)
    finally:
        if connection:
            connection.close()

url = []
for line in fileinput.input():
    urls.append(line)

crawl(urls)
