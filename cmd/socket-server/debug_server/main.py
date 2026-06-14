from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse

app = FastAPI()
app.mount("/static/",StaticFiles(directory="static"),name="static")

Base = ""

@app.get("/")
def root():
    return FileResponse("./static/pages/index.html")

@app.get("/viewer")
def viewer():
    return FileResponse("./static/pages/viewer.html")

@app.get("/table")
def nametable():
    return FileResponse("./static/pages/table.html")

@app.get("/main")
def main():
    return FileResponse("./static/pages/main.html")

@app.get("/socket")
def socket():
    return FileResponse("./static/pages/socket.html")

@app.get("/wasm")
def wasm():
    return FileResponse("./static/pages/wasm.html")