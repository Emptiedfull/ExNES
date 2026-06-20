from typing import Callable

from fastapi import FastAPI, Request
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint
from starlette.responses import Response

app = FastAPI()
app.mount("/static/",StaticFiles(directory="static"),name="static")

class CMW(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        response = await call_next(request)
        
        response.headers['Cross-Origin-Opener-Policy'] = 'same-origin'
        response.headers['Cross-Origin-Embedder-Policy'] = 'require-corp'
        
        return response

app.add_middleware(CMW)
Base = ""

@app.get("/")
def root():
    return FileResponse("./static/prototype.html")