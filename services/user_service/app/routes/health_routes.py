from fastapi import APIRouter
from app.utils.response import response


router = APIRouter()
@router.get("/health", tags=["health"])
async def health():
  return response(True, data={"status": "ok"}, message="healthy")