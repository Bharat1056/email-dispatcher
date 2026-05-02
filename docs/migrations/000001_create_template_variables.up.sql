CREATE TABLE template_variables (
    id UUID PRIMARY KEY,

    template_id UUID NOT NULL
        REFERENCES email_templates(id)
        ON DELETE CASCADE,

    key TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type in ('string', 'number', 'bool', 'date')),
    required BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW(),

    -- Prevent duplicate variables inside a template
    CONSTRAINT unique_variable_per_template 
        UNIQUE (template_id, key)
);

-- Fast lookup of variables for a template
CREATE INDEX idx_template_variables_template_id 
ON template_variables(template_id);

-- If frequently search by key
CREATE INDEX idx_template_variables_key 
ON template_variables(key);