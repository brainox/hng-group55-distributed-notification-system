#!/usr/bin/env python3
from typing import AsyncGenerator
from contextlib import asynccontextmanager

from sqlmodel import Session
from fastapi import FastAPI, status, Depends
from fastapi.responses import JSONResponse
from fastapi.background import BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware

from logger import Logger
from service.senders import FCMSender
from db.database import SessionLocal, get_db, init_db
from service.queue import RabbitMQQueue
from db.models import PushMessage, StatusEnum
from db.schema import PushMessageCreate, RootResponse

queue: RabbitMQQueue
background: BackgroundTasks


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:

    global queue

    Logger.info("Initiated Lifespan async function")

    queue = RabbitMQQueue(amqp_url="localhost:45672", queue_name="push_queue")

    init_db()

    yield

    queue.cleanup()
    Logger.info("Closing Lifespan async function")


app = FastAPI(
    title="Push Notification Service",
    version="1.0.0",
    description="A service to send push notifications via FCM.",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/", tags=["Root"], response_model=RootResponse)
async def read_root():
    Logger.info("Calling the root API")
    response = RootResponse(
        success=True,
        data=dict(
            title="Push Notification Service",
            version="1.0.0",
            description="A service to send push notifications via FCM.",
        ),
        error=None,
        message="RootAPI",
        meta=None,
    )

    return JSONResponse(content=response.model_dump(), status_code=status.HTTP_200_OK)


@app.get("/health", tags=["Health"], response_model=RootResponse)
async def health_check():
    Logger.info("Calling the health API")
    response = RootResponse(
        success=True,
        data={"status": "ok"},
        error=None,
        message="Health Check OK",
        meta=None,
    )

    return JSONResponse(content=response.model_dump(), status_code=status.HTTP_200_OK)


@app.post("/send", tags=["FCM"])
async def send_fcm_notification(
    payload: PushMessageCreate, db: Session = Depends(get_db)
):
    Logger.info("Received /send request with payload: %s", payload.model_dump())
    try:
        new_msg = PushMessage(
            id=payload.id,
            title=payload.title,
            body=payload.body,
            token=payload.token,
            status=StatusEnum.pending,
        )
        if not new_msg.token.strip():
            raise ValueError("Token is required")
        if not new_msg.title.strip():
            raise ValueError("Title is required")
        if not new_msg.body.strip():
            raise ValueError("Body is required")
        if not new_msg.id:
            raise ValueError("ID is required")
        if db is None:
            raise RuntimeError("Database session is not available")
        if new_msg is None:
            raise RuntimeError("Failed to create PushMessage instance")
        if db.get(PushMessage, new_msg.id):
            raise ValueError(f"PushMessage with ID {new_msg.id} already exists")
        db.add(new_msg)
        db.commit()
        db.refresh(new_msg)

        queue.push(
            {
                "provider": "fcm",
                "id": new_msg.id,
                "title": new_msg.title,
                "body": new_msg.body,
                "token": new_msg.token,
            }
        )
        response = RootResponse(
            success=True,
            data=new_msg,
            error=None,
            message="Health Check OK",
            meta=None,
        )

        return JSONResponse(
            content=response.model_dump(), status_code=status.HTTP_200_OK
        )
    except Exception as e:
        Logger.error("Error in /send: %s", str(e))
        response = RootResponse(
            success=False,
            data=None,
            error=str(e),
            message="Failed to send notification",
            meta=None,
        )
        return JSONResponse(
            content=response.model_dump(),
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        )


@app.get("/notifications", tags=["FCM"], response_model=RootResponse)
async def get_all_notifications():
    Logger.info("Received /notifications request")
    try:
        db = SessionLocal()
        msgs = db.query(PushMessage).all()
        response = RootResponse(
            success=True,
            data=msgs,
            error=None,
            message="Notifications retrieved successfully",
            meta=None,
        )
        return JSONResponse(
            content=response.model_dump(), status_code=status.HTTP_200_OK
        )
    except Exception as e:
        Logger.error("Error in /notifications: %s", str(e))
        response = RootResponse(
            success=False,
            data=None,
            error=str(e),
            message="Failed to retrieve notifications",
            meta=None,
        )
        return JSONResponse(
            content=response.model_dump(),
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        )


@app.get("/notifications/{notification_id}", tags=["FCM"], response_model=RootResponse)
async def get_single_notification(notification_id: int):
    Logger.info("Received /notifications/%s request", notification_id)
    try:
        db = SessionLocal()
        msg = db.get(PushMessage, notification_id)
        if not msg:
            response = RootResponse(
                success=False,
                data=None,
                error="Notification not found",
                message="Notification not found",
                meta=None,
            )
            return JSONResponse(
                content=response.model_dump(),
                status_code=status.HTTP_404_NOT_FOUND,
            )

        response = RootResponse(
            success=True,
            data=msg,
            error=None,
            message="Notification retrieved successfully",
            meta=None,
        )
        return JSONResponse(
            content=response.model_dump(), status_code=status.HTTP_200_OK
        )
    except Exception as e:
        Logger.error("Error in /notifications/%s: %s", notification_id, str(e))
        response = RootResponse(
            success=False,
            data=None,
            error=str(e),
            message="Failed to retrieve notification",
            meta=None,
        )
        return JSONResponse(
            content=response.model_dump(),
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        )
