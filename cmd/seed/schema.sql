CREATE TABLE "dag" (
  "id" INTEGER NOT NULL,
  "createdAt" integer,
  "updatedAt" integer,
  "status" integer NOT NULL,
  "parameters" text,
  "jobId" INTEGER,
  "dagTemplateId" INTEGER,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tasks_statuses_1" FOREIGN KEY ("status") REFERENCES "statuses" ("id"),
  CONSTRAINT "fk_tasks_jobs_1" FOREIGN KEY ("jobId") REFERENCES "jobs" ("id"),
  CONSTRAINT "fk_dag_dagTemplates_1" FOREIGN KEY ("dagTemplateId") REFERENCES "dagTemplates" ("id")
);

CREATE TABLE "dagTemplates" (
  "id" INTEGER NOT NULL,
  "name" TEXT,
  "version" integer,
  "jsonSchemaParameters" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "edges" (
  "id" INTEGER NOT NULL,
  "sourceTaskVerticeId" INTEGER,
  "destTaskVerticeId" INTEGER,
  "metadata" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "jobs" (
  "id" INTEGER NOT NULL,
  "createdAt" integer,
  "deletedAt" integer,
  "updatedAt" integer,
  "namespace" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "log" (
  "id" INTEGER NOT NULL,
  "dagId" INTEGER NOT NULL,
  "createdAt" integer,
  "logLevel" integer,
  "message" TEXT,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_log_dag_1" FOREIGN KEY ("id") REFERENCES "dag" ("id")
);

CREATE TABLE "metadata" (
  "dagId" INTEGER NOT NULL,
  "metadata" TEXT,
  PRIMARY KEY ("dagId"),
  CONSTRAINT "fk_metadata_dag_1" FOREIGN KEY ("dagId") REFERENCES "dag" ("id")
);

CREATE TABLE "statuses" (
  "id" INTEGER NOT NULL,
  "label" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "table_1" (

);

CREATE TABLE "templateEdges" (
  "id" INTEGER NOT NULL,
  "sourceTaskVerticeId" INTEGER,
  "destTaskVerticeId" INTEGER,
  "metadata" TEXT,
  PRIMARY KEY ("id")
);

CREATE TABLE "templateVertices" (
  "id" INTEGER NOT NULL,
  "json" TEXT,
  "coordinates" TEXT,
  "typeId" INTEGER NOT NULL,
  "templateId" INTEGER NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_taskVertices_taskEdges_1_copy_1" FOREIGN KEY ("id") REFERENCES "templateEdges" ("sourceTaskVerticeId"),
  CONSTRAINT "fk_taskVertices_taskEdges_2_copy_1" FOREIGN KEY ("id") REFERENCES "templateEdges" ("destTaskVerticeId"),
  CONSTRAINT "fk_vertices_verticesTypes_1_copy_1" FOREIGN KEY ("typeId") REFERENCES "verticesTypes" ("id"),
  CONSTRAINT "fk_templateVertices_dagTemplates_1" FOREIGN KEY ("templateId") REFERENCES "dagTemplates" ("id")
);

CREATE TABLE "vertices" (
  "id" INTEGER NOT NULL,
  "status" integer NOT NULL,
  "json" TEXT,
  "coordinates" TEXT,
  "typeId" INTEGER NOT NULL,
  "dagId" INTEGER NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_taskNodes_statuses_1" FOREIGN KEY ("status") REFERENCES "statuses" ("id"),
  CONSTRAINT "fk_taskVertices_taskEdges_1" FOREIGN KEY ("id") REFERENCES "edges" ("sourceTaskVerticeId"),
  CONSTRAINT "fk_taskVertices_taskEdges_2" FOREIGN KEY ("id") REFERENCES "edges" ("destTaskVerticeId"),
  CONSTRAINT "fk_vertices_verticesTypes_1" FOREIGN KEY ("typeId") REFERENCES "verticesTypes" ("id"),
  CONSTRAINT "fk_vertices_dag_1" FOREIGN KEY ("dagId") REFERENCES "dag" ("id")
);

CREATE TABLE "verticesTypes" (
  "id" INTEGER NOT NULL,
  "label" TEXT,
  PRIMARY KEY ("id")
);

