from fastapi import FastAPI
from app.routes import user_router, health_router,auth_router, preferences_router
from app.core.db import init_db, init_redis

app = FastAPI(title="user_service")


@app.on_event("startup")
async def startup():
  await init_db()
  await init_redis()





app.include_router(user_router)
app.include_router(health_router)
app.include_router(auth_router)
app.include_router(preferences_router)