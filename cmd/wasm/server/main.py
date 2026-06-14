from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse

app = FastAPI()
app.mount("/static/",StaticFiles(directory="static"),name="static")

Base = ""

@app.get("/")
def root():
    return FileResponse("./static/index.html")



