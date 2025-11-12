import pytest
from httpx import AsyncClient
from app.main import app

@pytest.mark.asyncio
async def test_health_endpoint():
    async with AsyncClient(app=app, base_url="http://test") as ac:
        response = await ac.get("/health")
        assert response.status_code == 200
        data = response.json()
        assert data["success"] is True

@pytest.mark.asyncio
async def test_register_and_login():
    async with AsyncClient(app=app, base_url="http://test") as ac:
        # Register user
        payload = {
            "name": "Test User",
            "email": "test@example.com",
            "password": "password123"
        }
        response = await ac.post("/v1/auth/register", json=payload)
        assert response.status_code == 200
        data = response.json()
        assert data["success"] is True
        user_id = data["data"]["user_id"]

        # Login
        login_payload = {
            "email": "test@example.com",
            "password": "password123"
        }
        response = await ac.post("/v1/auth/login", json=login_payload)
        assert response.status_code == 200
        token_data = response.json()["data"]
        assert "access_token" in token_data

        # Get preferences (should fail initially)
        response = await ac.get(f"/v1/preferences/{user_id}")
        assert response.status_code == 404
