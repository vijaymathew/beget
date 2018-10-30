from flask import Flask, request, jsonify
app = Flask('httprepo')

@app.route('/doc', methods=['POST'])
def post_doc():
    content = request.json
    print(content)
    return 'OK'

app.run(host='localhost', port=5050)
