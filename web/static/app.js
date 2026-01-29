/**
 * Aliasly Web UI - Frontend JavaScript
 *
 * This file handles all the frontend logic for the web configuration interface:
 * - Fetching and displaying aliases
 * - Creating, editing, and deleting aliases
 * - Form handling and validation
 */

// ============================================
// API Functions
// ============================================

/**
 * Fetches all aliases from the server.
 * @returns {Promise<Array>} Array of alias objects
 */
async function fetchAliases() {
    const response = await fetch('/api/aliases');
    const result = await response.json();

    if (!result.success) {
        throw new Error(result.error || 'Failed to fetch aliases');
    }

    return result.data || [];
}

/**
 * Creates a new alias on the server.
 * @param {Object} alias - The alias object to create
 * @returns {Promise<Object>} The created alias
 */
async function createAlias(alias) {
    const response = await fetch('/api/aliases', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(alias)
    });

    const result = await response.json();

    if (!result.success) {
        throw new Error(result.error || 'Failed to create alias');
    }

    return result.data;
}

/**
 * Updates an existing alias on the server.
 * @param {string} name - The name of the alias to update
 * @param {Object} alias - The updated alias object
 * @returns {Promise<Object>} The updated alias
 */
async function updateAlias(name, alias) {
    const response = await fetch(`/api/aliases/${encodeURIComponent(name)}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(alias)
    });

    const result = await response.json();

    if (!result.success) {
        throw new Error(result.error || 'Failed to update alias');
    }

    return result.data;
}

/**
 * Deletes an alias from the server.
 * @param {string} name - The name of the alias to delete
 */
async function deleteAlias(name) {
    const response = await fetch(`/api/aliases/${encodeURIComponent(name)}`, {
        method: 'DELETE'
    });

    const result = await response.json();

    if (!result.success) {
        throw new Error(result.error || 'Failed to delete alias');
    }
}

// ============================================
// UI Rendering (Using safe DOM methods)
// ============================================

/**
 * Renders the list of aliases to the page.
 * @param {Array} aliases - Array of alias objects
 */
function renderAliases(aliases) {
    const container = document.getElementById('aliasList');
    container.textContent = ''; // Clear existing content safely

    // Handle empty state
    if (!aliases || aliases.length === 0) {
        const emptyDiv = document.createElement('div');
        emptyDiv.className = 'empty-state';

        const h3 = document.createElement('h3');
        h3.textContent = 'No aliases yet';
        emptyDiv.appendChild(h3);

        const p = document.createElement('p');
        p.textContent = 'Click "Add New Alias" to create your first alias.';
        emptyDiv.appendChild(p);

        container.appendChild(emptyDiv);
        return;
    }

    // Render each alias as a card
    for (const alias of aliases) {
        container.appendChild(createAliasCard(alias));
    }
}

/**
 * Creates a single alias card element.
 * @param {Object} alias - The alias object
 * @returns {HTMLElement} The card element
 */
function createAliasCard(alias) {
    const card = document.createElement('div');
    card.className = 'alias-card';
    card.dataset.alias = alias.name;

    // Header
    const header = document.createElement('div');
    header.className = 'alias-header';

    const nameSpan = document.createElement('span');
    nameSpan.className = 'alias-name';
    nameSpan.textContent = alias.name;
    header.appendChild(nameSpan);

    // Action buttons
    const actions = document.createElement('div');
    actions.className = 'alias-actions';

    const editBtn = document.createElement('button');
    editBtn.className = 'btn-icon';
    editBtn.title = 'Edit';
    editBtn.onclick = () => editAlias(alias.name);
    editBtn.innerHTML = '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>';
    actions.appendChild(editBtn);

    const deleteBtn = document.createElement('button');
    deleteBtn.className = 'btn-icon danger';
    deleteBtn.title = 'Delete';
    deleteBtn.onclick = () => confirmDelete(alias.name);
    deleteBtn.innerHTML = '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>';
    actions.appendChild(deleteBtn);

    header.appendChild(actions);
    card.appendChild(header);

    // Description
    if (alias.description) {
        const desc = document.createElement('div');
        desc.className = 'alias-description';
        desc.textContent = alias.description;
        card.appendChild(desc);
    }

    // Command
    const cmd = document.createElement('div');
    cmd.className = 'alias-command';
    cmd.textContent = alias.command;
    card.appendChild(cmd);

    // Parameters
    if (alias.params && alias.params.length > 0) {
        const paramsDiv = document.createElement('div');
        paramsDiv.className = 'alias-params';
        paramsDiv.appendChild(document.createTextNode('Params: '));

        for (const p of alias.params) {
            const span = document.createElement('span');
            if (p.required) {
                span.className = 'required';
            }
            let text = p.name;
            if (p.required) text += '*';
            if (p.default) text += ` = ${p.default}`;
            span.textContent = text;
            paramsDiv.appendChild(span);
        }

        card.appendChild(paramsDiv);
    }

    // Usage
    const usageDiv = document.createElement('div');
    usageDiv.className = 'alias-usage';
    usageDiv.appendChild(document.createTextNode('Usage: '));

    const usageCode = document.createElement('code');
    usageCode.textContent = 'al ' + buildUsageString(alias);
    usageDiv.appendChild(usageCode);

    card.appendChild(usageDiv);

    return card;
}

/**
 * Builds a usage string for an alias.
 * @param {Object} alias - The alias object
 * @returns {string} Usage string like "gc <message>" or "gp [branch]"
 */
function buildUsageString(alias) {
    let usage = alias.name;

    if (alias.params) {
        for (const p of alias.params) {
            if (p.required) {
                usage += ` <${p.name}>`;
            } else {
                usage += ` [${p.name}]`;
            }
        }
    }

    return usage;
}

// ============================================
// Modal Functions
// ============================================

// Track which alias we're editing (null = creating new)
let editingAlias = null;

/**
 * Opens the modal for adding a new alias.
 */
function openAddModal() {
    editingAlias = null;
    document.getElementById('modalTitle').textContent = 'Add New Alias';
    document.getElementById('aliasForm').reset();
    document.getElementById('aliasName').disabled = false;
    document.getElementById('paramsContainer').textContent = '';
    updatePreview();
    document.getElementById('modal').classList.remove('hidden');
}

/**
 * Opens the modal for editing an existing alias.
 * @param {string} name - The name of the alias to edit
 */
async function editAlias(name) {
    try {
        const aliases = await fetchAliases();
        const alias = aliases.find(a => a.name === name);

        if (!alias) {
            alert('Alias not found');
            return;
        }

        editingAlias = alias;
        document.getElementById('modalTitle').textContent = 'Edit Alias';

        // Populate form
        document.getElementById('aliasName').value = alias.name;
        document.getElementById('aliasName').disabled = true; // Can't change name
        document.getElementById('aliasCommand').value = alias.command;
        document.getElementById('aliasDescription').value = alias.description || '';

        // Populate params
        const paramsContainer = document.getElementById('paramsContainer');
        paramsContainer.textContent = '';

        if (alias.params) {
            for (const p of alias.params) {
                addParamField(p.name, p.description, p.required, p.default);
            }
        }

        updatePreview();
        document.getElementById('modal').classList.remove('hidden');
    } catch (error) {
        alert('Error loading alias: ' + error.message);
    }
}

/**
 * Closes the add/edit modal.
 */
function closeModal() {
    document.getElementById('modal').classList.add('hidden');
    editingAlias = null;
}

/**
 * Opens the delete confirmation modal.
 * @param {string} name - The name of the alias to delete
 */
function confirmDelete(name) {
    document.getElementById('deleteAliasName').textContent = name;
    document.getElementById('confirmDeleteBtn').onclick = () => performDelete(name);
    document.getElementById('deleteModal').classList.remove('hidden');
}

/**
 * Closes the delete confirmation modal.
 */
function closeDeleteModal() {
    document.getElementById('deleteModal').classList.add('hidden');
}

/**
 * Performs the actual deletion of an alias.
 * @param {string} name - The name of the alias to delete
 */
async function performDelete(name) {
    try {
        await deleteAlias(name);
        closeDeleteModal();
        await loadAliases();
    } catch (error) {
        alert('Error deleting alias: ' + error.message);
    }
}

// ============================================
// Parameter Fields
// ============================================

/**
 * Adds a parameter field to the form.
 * @param {string} name - Parameter name (optional, for editing)
 * @param {string} description - Parameter description (optional)
 * @param {boolean} required - Whether the parameter is required
 * @param {string} defaultValue - Default value (optional)
 */
function addParamField(name = '', description = '', required = true, defaultValue = '') {
    const container = document.getElementById('paramsContainer');

    const div = document.createElement('div');
    div.className = 'param-field';

    // Name input
    const nameInput = document.createElement('input');
    nameInput.type = 'text';
    nameInput.placeholder = 'Name';
    nameInput.value = name;
    nameInput.className = 'param-name';
    nameInput.onchange = updatePreview;
    div.appendChild(nameInput);

    // Description input
    const descInput = document.createElement('input');
    descInput.type = 'text';
    descInput.placeholder = 'Description';
    descInput.value = description;
    descInput.className = 'param-desc';
    div.appendChild(descInput);

    // Required select
    const reqSelect = document.createElement('select');
    reqSelect.className = 'param-required';
    reqSelect.onchange = updatePreview;

    const optRequired = document.createElement('option');
    optRequired.value = 'true';
    optRequired.textContent = 'Required';
    if (required) optRequired.selected = true;
    reqSelect.appendChild(optRequired);

    const optOptional = document.createElement('option');
    optOptional.value = 'false';
    optOptional.textContent = 'Optional';
    if (!required) optOptional.selected = true;
    reqSelect.appendChild(optOptional);

    div.appendChild(reqSelect);

    // Remove button
    const removeBtn = document.createElement('button');
    removeBtn.type = 'button';
    removeBtn.className = 'btn-icon danger';
    removeBtn.title = 'Remove';
    removeBtn.onclick = () => removeParamField(removeBtn);
    removeBtn.innerHTML = '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>';
    div.appendChild(removeBtn);

    container.appendChild(div);
    updatePreview();
}

/**
 * Removes a parameter field from the form.
 * @param {HTMLElement} button - The remove button that was clicked
 */
function removeParamField(button) {
    button.parentElement.remove();
    updatePreview();
}

/**
 * Collects all parameters from the form.
 * @returns {Array} Array of parameter objects
 */
function collectParams() {
    const container = document.getElementById('paramsContainer');
    const params = [];

    for (const field of container.children) {
        const name = field.querySelector('.param-name').value.trim();
        if (!name) continue;

        const param = {
            name: name,
            description: field.querySelector('.param-desc').value.trim(),
            required: field.querySelector('.param-required').value === 'true',
            default: ''
        };

        const defaultInput = field.querySelector('.param-default');
        if (defaultInput) {
            param.default = defaultInput.value.trim();
        }

        params.push(param);
    }

    return params;
}

// ============================================
// Preview
// ============================================

/**
 * Updates the usage preview based on current form values.
 */
function updatePreview() {
    const name = document.getElementById('aliasName').value.trim() || 'alias';
    const params = collectParams();

    let usage = `al ${name}`;
    for (const p of params) {
        if (p.required) {
            usage += ` <${p.name}>`;
        } else {
            usage += ` [${p.name}]`;
        }
    }

    document.getElementById('usagePreview').textContent = usage;
}

// ============================================
// Form Submission
// ============================================

/**
 * Handles form submission for creating/updating aliases.
 * @param {Event} event - The form submit event
 */
async function handleSubmit(event) {
    event.preventDefault();

    const alias = {
        name: document.getElementById('aliasName').value.trim(),
        command: document.getElementById('aliasCommand').value.trim(),
        description: document.getElementById('aliasDescription').value.trim(),
        params: collectParams()
    };

    // Remove empty params array
    if (alias.params.length === 0) {
        delete alias.params;
    }

    try {
        if (editingAlias) {
            await updateAlias(editingAlias.name, alias);
        } else {
            await createAlias(alias);
        }

        closeModal();
        await loadAliases();
    } catch (error) {
        alert('Error saving alias: ' + error.message);
    }
}

// ============================================
// Utility Functions
// ============================================

/**
 * Loads and displays all aliases.
 */
async function loadAliases() {
    const container = document.getElementById('aliasList');

    try {
        const aliases = await fetchAliases();
        allAliases = aliases; // Store for search
        renderAliases(aliases);
    } catch (error) {
        container.textContent = '';

        const errorDiv = document.createElement('div');
        errorDiv.className = 'empty-state';

        const h3 = document.createElement('h3');
        h3.textContent = 'Error loading aliases';
        errorDiv.appendChild(h3);

        const p = document.createElement('p');
        p.textContent = error.message;
        errorDiv.appendChild(p);

        container.appendChild(errorDiv);
    }
}

// ============================================
// Event Listeners
// ============================================

// ============================================
// Event Listeners
// ============================================

// When the page loads, fetch and display aliases
document.addEventListener('DOMContentLoaded', () => {
    // Initialize theme
    initTheme();

    loadAliases();

    // Set up event listeners
    document.getElementById('addAliasBtn').addEventListener('click', () => openAddModal());
    document.getElementById('aliasForm').addEventListener('submit', handleSubmit);
    document.getElementById('themeToggle').addEventListener('click', toggleTheme);
    document.getElementById('searchInput').addEventListener('input', handleSearch);
    document.getElementById('exportBtn').addEventListener('click', exportConfig);
    document.getElementById('importBtn').addEventListener('click', () => document.getElementById('importFileInput').click());
    document.getElementById('importFileInput').addEventListener('change', handleImport);

    // Update preview when command changes (to detect {{params}})
    document.getElementById('aliasCommand').addEventListener('input', () => {
        // Auto-detect parameters from command
        const command = document.getElementById('aliasCommand').value;
        const matches = command.match(/\{\{(\w+)\}\}/g) || [];
        const paramNames = matches.map(m => m.slice(2, -2));

        // Get existing param names
        const existingParams = collectParams().map(p => p.name);

        // Add fields for new params
        for (const name of paramNames) {
            if (!existingParams.includes(name)) {
                addParamField(name, '', true, '');
            }
        }

        updatePreview();
    });

    document.getElementById('aliasName').addEventListener('input', updatePreview);

    // Close modal when clicking outside
    document.getElementById('modal').addEventListener('click', (e) => {
        if (e.target.id === 'modal') closeModal();
    });

    document.getElementById('deleteModal').addEventListener('click', (e) => {
        if (e.target.id === 'deleteModal') closeDeleteModal();
    });

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            closeModal();
            closeDeleteModal();
        }
    });
});

// ============================================
// Theme Handling
// ============================================

function initTheme() {
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

    if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
        setTheme('dark');
    } else {
        setTheme('light');
    }
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    setTheme(newTheme);
}

function setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);

    // Update icon visibility
    const sunIcon = document.querySelector('.sun-icon');
    const moonIcon = document.querySelector('.moon-icon');

    if (theme === 'dark') {
        sunIcon.style.display = 'none';
        moonIcon.style.display = 'block';
    } else {
        sunIcon.style.display = 'block';
        moonIcon.style.display = 'none';
    }
}

// ============================================
// Search Handling
// ============================================

let allAliases = [];

function handleSearch(e) {
    const query = e.target.value.toLowerCase();
    const filtered = allAliases.filter(alias =>
        alias.name.toLowerCase().includes(query) ||
        alias.command.toLowerCase().includes(query) ||
        (alias.description && alias.description.toLowerCase().includes(query))
    );
    renderAliases(filtered);
}

// ============================================
// Import/Export Functions
// ============================================

/**
 * Exports the config file by triggering a download.
 */
function exportConfig() {
    // Create a link to download the config
    const link = document.createElement('a');
    link.href = '/api/config/export';
    link.download = 'aliasly-config.yaml';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

/**
 * Handles file import when user selects a file.
 * @param {Event} event - The file input change event
 */
async function handleImport(event) {
    const file = event.target.files[0];
    if (!file) return;

    // Confirm import
    const confirmImport = confirm(
        `Import "${file.name}"?\n\nNew aliases will be added to your existing configuration.\nAliases with the same name will be skipped.`
    );

    if (!confirmImport) {
        event.target.value = ''; // Reset file input
        return;
    }

    // Create form data
    const formData = new FormData();
    formData.append('config', file);

    try {
        const response = await fetch('/api/config/import', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (!result.success) {
            throw new Error(result.error || 'Failed to import config');
        }

        // Update aliases from result
        const importResult = result.data;
        allAliases = importResult.aliases || [];
        renderAliases(allAliases);

        // Show result message
        let message = `Import complete!\n\nAdded: ${importResult.added} alias(es)`;
        if (importResult.skipped > 0) {
            message += `\nSkipped: ${importResult.skipped} (already exist)`;
        }
        alert(message);
    } catch (error) {
        alert('Error importing config: ' + error.message);
    }

    // Reset file input so same file can be selected again
    event.target.value = '';
}
