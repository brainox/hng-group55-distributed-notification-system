#!/usr/bin/env python3
from typing import AsyncGenerator
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.background import BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware

from service.worker import PushWorker
from service.senders import FCMSender
from service.queue import RabbitMQQueue
from logger import Logger

sender: FCMSender
worker: PushWorker
queue: RabbitMQQueue
background: BackgroundTasks


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    
    Logger.info("Initiated Lifespan async function")

    global sender, worker, queue, background

    sender = FCMSender(
        credential_path="distributed-systems-3349d-firebase-adminsdk-fbsvc-b84f851127.json"
    )
    queue = RabbitMQQueue(amqp_url="localhost:45672", queue_name="push_queue")

    worker = PushWorker(queue, sender)

    background = BackgroundTasks()
    background.add_task(worker.run_forever, poll_timeout=5)

    yield
    
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


@app.get("/", tags=["Root"])
async def read_root():
    return


@app.get("/health", tags=["Health"])
async def health_check():
    
    return {"status": "ok"}


@app.post("/send", tags=["FCM"])
async def send_fcm_notification():
    pass
