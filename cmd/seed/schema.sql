CREATE TABLE "jobs" (
  "id" INTEGER NOT NULL,
  "createdAt" integer,
  "deletedAt" integer,
  "updatedAt" integer,
  "namespace" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "statuses" (
  "id" INTEGER NOT NULL,
  "label" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "taskEdges" (
  "id" INTEGER NOT NULL,
  "sourceTaskVerticeId" INTEGER,
  "destTaskVerticeId" INTEGER,
  PRIMARY KEY ("id")
);

CREATE TABLE "tasks" (
  "id" INTEGER NOT NULL,
  "status" integer NOT NULL,
  "jobId" INTEGER,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tasks_statuses_1" FOREIGN KEY ("status") REFERENCES "statuses" ("id"),
  CONSTRAINT "fk_tasks_jobs_1" FOREIGN KEY ("jobId") REFERENCES "jobs" ("id")
);

CREATE TABLE "taskVertices" (
  "id" INTEGER NOT NULL,
  "status" integer NOT NULL,
  "json" TEXT,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_taskNodes_statuses_1" FOREIGN KEY ("status") REFERENCES "statuses" ("id"),
  CONSTRAINT "fk_taskVertices_taskEdges_1" FOREIGN KEY ("id") REFERENCES "taskEdges" ("sourceTaskVerticeId"),
  CONSTRAINT "fk_taskVertices_taskEdges_2" FOREIGN KEY ("id") REFERENCES "taskEdges" ("destTaskVerticeId")
);

