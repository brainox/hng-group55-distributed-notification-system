from fastapi import FastAPI
from app.routes import user_router, health_router, auth_router
from app.core.db import init_db
import asyncio


app = FastAPI(title="user_service")


@app.on_event("startup")
async def startup():
    # Initialize DB
    await init_db()
   

  



# Routers
app.include_router(user_router)
app.include_router(health_router)
app.include_router(auth_router)

