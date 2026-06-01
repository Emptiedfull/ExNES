from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse

app = FastAPI()
app.mount("/static/",StaticFiles(directory="static"),name="static")

@app.get("/")
def root():
    return FileResponse("C:/Users/user/ExNES/debug_server/static/pages/index.html")

@app.get("/viewer")
def viewer():
    return FileResponse("C:/Users/user/ExNES/debug_server/static/pages/viewer.html")