import asyncio
from httpx import AsyncClient
from app.main import app


async def test_health():
    async with AsyncClient(app=app, base_url="http://test") as ac:
          res = await ac.get("/health")
          assert res.status_code == 200
          body = res.json()
          assert body["success"] is True