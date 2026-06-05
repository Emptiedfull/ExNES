from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse

app = FastAPI()
app.mount("/static/",StaticFiles(directory="static"),name="static")

@app.get("/")
def root():
    return FileResponse("C:/Users/user/ExNES/emulatorCore/debug_server/static/pages/index.html")

@app.get("/viewer")
def viewer():
    return FileResponse("C:/Users/user/ExNES/emulatorCore/debug_server/static/pages/viewer.html")

@app.get("/nametable")
def nametable():
    return FileResponse("C:/Users/user/ExNES/debug_server/static/pages/table.html")

@app.get("/main")
def main():
    return FileResponse("C:/Users/user/ExNES/emulatorCore/debug_server/static/pages/main.html")

@app.get("/socket")
def socket():
    return FileResponse("C:/Users/user/ExNES/emulatorCore/debug_server/static/pages/socket.html")