import http.client
import json
import sys
import ssl
import getopt
import fileinput
import traceback

def print_usage():
    print('conroller.py <options>\n')
    print('options: \n')
    print('-h                      print this help and quit\n')
    print('-c --crawlerhost=URL    url to reach the crawler service\n')
    print('-u --urls=FILENAME      name of file with URLs to crawl, one on each line\n')
    print('-t --repotype=TYPE      the type of repository (file, simpleHTTP etc)\n')
    print('-r --repoconfig=CONFIG  configuration for the repository\n')

crawler_url = ''
crawl_seed_file = ''
crawl_repo_type = 'file'
crawl_repo_config = '.'
opts = []
args = []

try:
    opts, args = getopt.getopt(sys.argv[1:], "hc:u:r:t:",["crawlerhost=","urls=", "repoconfig=", "repotype="])
except getopt.GetoptError:
    print_usage()
    sys.exit(2)

for opt, arg in opts:
    if opt == '-h':
        print_usage()
        sys.exit()
    elif opt in ("-c", "--crawlerhost"):
        crawler_url = arg
    elif opt in ("-u", "--urls"):
        crawl_seed_file = arg
    elif opt in ("-t", "--repotype"):
        crawl_repo_type = arg
    elif opt in ("-r", "--repoconfig"):
        crawl_repo_config = arg

crawler_host = crawler_url[crawler_url.find('://')+3:]

def urlfile(url):
    return url[url.rfind("/")+1:]

def getconnection():
    try:
        if crawler_url.startswith('https'):
            ctx = ssl._create_unverified_context()
            ctx.check_hostname = False
            ctx.verify_mode = ssl.CERT_NONE
            return http.client.HTTPSConnection(crawler_host, context = ctx)
        else:
            return http.client.HTTPConnection(crawler_host)
    except:
        traceback.print_exc(file=sys.stdout)

headers = {'Content-type': 'application/json'}

def crawl(urls):
    resources = {}
    for url in urls:
        k = urlfile(url)
        resources[k] = url

    request = {"repository": crawl_repo_type,
               "repositoryConfig": crawl_repo_config,
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
        traceback.print_exc(file=sys.stdout)
    finally:
        if connection:
            connection.close()

urls = []
for line in fileinput.input(crawl_seed_file):
    urls.append(line.strip())

crawl(urls)
