from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.future import select
from app.models.preference import Preference
from app.schemas.preference_schema import PreferenceCreate, PreferenceOut
from app.utils.response import response
from app.core.db import async_session
import uuid

router = APIRouter(prefix="/v1/preferences", tags=["preferences"])

async def get_db():
    async with async_session() as session:
        yield session


@router.post("/{user_id}", response_model=dict)
async def create_or_update_preferences(user_id: uuid.UUID, payload: PreferenceCreate, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Preference).where(Preference.user_id == user_id))
    existing = result.scalar_one_or_none()

    if existing:
        existing.email = payload.email
        existing.push = payload.push
        await db.commit()
        await db.refresh(existing)
        return response(True, data=PreferenceOut.from_orm(existing), message="Preferences updated")

    pref = Preference(user_id=user_id, channel='default', template_type='user_pref', enabled=True)
    pref.email = payload.email
    pref.push = payload.push
    db.add(pref)
    await db.commit()
    await db.refresh(pref)
    return response(True, data=PreferenceOut.from_orm(pref), message="Preferences created")


@router.get("/{user_id}", response_model=dict)
async def get_preferences(user_id: uuid.UUID, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Preference).where(Preference.user_id == user_id))
    pref = result.scalar_one_or_none()
    if not pref:
        raise HTTPException(status_code=404, detail="User preferences not found")
    return response(True, data=PreferenceOut.from_orm(pref), message="User preferences retrieved")
