CREATE TABLE IF NOT EXISTS email_templates (
    template_id UUID PRIMARY KEY,

    name TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'draft' CHECK(status in ('draft', 'active', 'inactive')), -- draft | active | inactive

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    created_by UUID NOT NULL 
        REFERENCES users(id)
    
    CONSTRAINT unique_template_per_user
        UNIQUE (created_by, name)
);

-- Fetch templates by user
CREATE INDEX idx_email_templates_created_by ON email_templates(created_by);

-- Filter by status (active/draft)
CREATE INDEX idx_email_templates_status ON email_templates(status);

-- Sorting / pagination
CREATE INDEX idx_email_templates_created_at ON email_templates(created_at DESC);