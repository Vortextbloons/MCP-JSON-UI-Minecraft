# Minecraft Bedrock JSON UI MCP Server

MCP server (Go, stdio) that helps AI agents author valid Minecraft Bedrock JSON UI. It exposes curated docs, UI generators, a structural validator, and binding/inheritance analysis tools centered on the **Server Form Hijack** pattern.

Summary: Bedrock has no first-class custom UI API. This server teaches the workaround — trigger a `server-ui` Server Form from a behavior pack, override `server_form.json` in a resource pack to inject your own layout — and then generates / validates that JSON so it works in-game.

- Module: `bedrock-json-ui-mcp` (`go.mod`, Go 1.26.2)
- Entry: `main.go` -> `internal/app/app.go:11` (`app.Run()` blocks on SIGINT/SIGTERM)
- Transport: stdio via `github.com/metoro-io/mcp-golang` (`internal/mcpserver/server.go:56`)
- Server name/version: `minecraft-bedrock-json-ui` / `2.0.0`
- Validation: `github.com/santhosh-tekuri/jsonschema/v5` against embedded schema (`internal/bedrockui/validator.go`, `schema.go`)

## Server Form Hijack (core pattern)

1. `resource_pack/ui/my_custom_ui.json` — your layout under a custom `namespace` (never `server_form`).
2. `resource_pack/ui/_ui_defs.json` — register your file + `server_form.json`.
3. `resource_pack/ui/server_form.json` — override vanilla screen (`namespace: server_form`), route `server_form_factory` / `main_screen_content` to your layout, use `#title_text` view binding for title-based routing.
4. `behavior_pack/scripts/main.js` — `ActionFormData().title(...).show(player)` via `@minecraft/server-ui`. Labels/values/options live in the script, not the JSON.

Full guide is served as resource `bedrock://docs/custom_json_ui_design` (`internal/bedrockui/reference.go:64`).

## Tools

| Tool | File | What it does |
|------|------|--------------|
| `validate_ui_contract` | `internal/mcpserver/handlers.go:28`, `internal/bedrockui/validator.go:40` | Validates a stringified JSON UI doc against the local partial schema. Returns success or numbered errors. Always prepends implementation hints (root namespace required, `controls` semantics, view vs collection-details bindings, root-level `anims`, `factory` usage). Partial check only — does not verify templates, external namespaces, or in-game behavior. Input: `ui_json: string`. |
| `generate_action_form` | `internal/mcpserver/handlers.go:133`, `internal/bedrockui/generators.go:63` | Generates action-form layout + install guidance (`server_form.json` companion patch, `_ui_defs.json` step). Inputs: `namespace` (default `custom_ui`), `variant`: `default`/`fullscreen`/`grid`/`tabs`/`light`, `icon_support: bool`, `include_scroll: bool`. Uses factory `buttons` + `form_buttons` collection, Leaf/common inherited buttons. |
| `generate_modal_form` | `internal/mcpserver/handlers.go:150`, `internal/bedrockui/generators.go:91` | Generates modal-form template inheriting vanilla `server_form.custom_*` widgets. Inputs: `namespace`, `controls[]` (`type`: `label`/`toggle`/`slider`/`step_slider`/`dropdown`/`input`/`header`/`divider`/`multiselect`, plus `label`, `placeholder`, `default`, `min`/`max`, `step`, `options`). Emits `control_ids` factory mapping. Field data still comes from `ModalFormData` in BP script. |
| `generate_tabbed_form` | `internal/mcpserver/handlers.go:165`, `internal/bedrockui/generators.go:144` | Generates tabbed layout sketch. Input: `namespace`, `tab_names: string[]`. Layout only — markers do not switch tabs, script-driven page selection required. |
| `generate_grid_form` | `internal/mcpserver/handlers.go:182`, `internal/bedrockui/generators.go:254` | Generates grid action form. Inputs: `namespace`, `columns` (default 3), `rows` (default 4), `icon_support`. Uses `grid` + `grid_item_template` + `buttons` factory on `form_buttons` collection. |
| `cherry_theme_pack` | `internal/mcpserver/handlers.go:207`, `internal/bedrockui/generators.go:314` | Packs/unpacks Cherry UI / Nut UI theme strings (Leaf Essentials `§f§0§0` fixed-width format). `mode: pack` inputs: `title` (100ch), `header`/`border`/`h_button`/`de_button`/`di_button`/`ho_button`/`paperdoll` (50ch each). `mode: unpack` input: `packed_string`. |
| `analyze_bindings` | `internal/mcpserver/handlers.go:232`, `internal/bedrockui/analyzer.go:32` | Walks a doc, extracts all `bindings[]`. Returns `total_bindings`, `by_type`, `by_collection`, per-binding `path`/`binding_type`/`binding_name`/`override`/`collection`/`source_property_name`/`target_property_name`/`binding_condition`, plus warnings (e.g. missing `binding_name` except `collection_details`/`view`/`none`). Input: `ui_json`. |
| `expand_inheritance` | `internal/mcpserver/handlers.go:248`, `internal/bedrockui/analyzer.go:158` | Expands static local `name@namespace.control` inheritance into resolved tree (`type`, merged `properties`, `children`). Warns on external namespaces, cycles, unevaluated `modifications`. Does not evaluate variables. Input: `ui_json`. |

Generated-form tools return two content blocks: raw layout JSON + install text (save under `resource_pack/ui`, register in `ui/_ui_defs.json`, merge companion `server_form` override).

## Resources

| URI | Name | Contents |
|-----|------|----------|
| `bedrock://docs/ui_reference` | `ui_reference` | Quick curated reference: element types, layout/binding props, hardcoded `#bindings`. `internal/bedrockui/reference.go:43`. |
| `bedrock://docs/json_ui_full_reference` | `json_ui_full_reference` | Full restructured JSON UI docs (embedded `internal/bedrockui/full_reference.txt`). |
| `bedrock://docs/custom_json_ui_design` | `custom_json_ui_design` | Server Form Hijack implementation guide. |
| `bedrock://docs/leaf_button_templates` | `leaf_button_templates` | Leaf `light_text_button` / `light_content_button` / `light_glyph_button`, StarLib2 styles (`leaf_button_templates.txt`). |
| `bedrock://docs/leaf_dialog_layouts` | `leaf_dialog_layouts` | `main_panel_no_buttons` / `_one_button` / `_two_buttons`, backgrounds (`leaf_dialog_layouts.txt`). |
| `bedrock://docs/binding_patterns` | `binding_patterns` | Collection / view / global binding examples, string-op syntax, factory wiring (`binding_patterns.txt`). |
| `bedrock://docs/custom_form_controls` | `custom_form_controls` | Modal control definitions + binding names (`custom_form_controls.txt`). |

Recommended agent flow: read quick + full + design resources -> draft JSON -> `validate_ui_contract` -> `analyze_bindings` / `expand_inheritance` if bindings/inheritance are complex -> fix -> validate again -> test in-game with content logs.

## Build / Run

```powershell
go build -o Bedrock-JSON-UI-MCP.exe .
.\Bedrock-JSON-UI-MCP.exe
```

Runs on stdio; block stays alive until Ctrl+C. Test with:

```powershell
go test ./...
```

MCP client config (example):

```json
{
  "mcpServers": {
    "minecraft-bedrock-json-ui": {
      "command": "C:\\path\\to\\Bedrock-JSON-UI-MCP.exe"
    }
  }
}
```

## Limitations

- Validator is structural/partial. Passing it does not guarantee engine compatibility. Verify inherited templates against your target game's vanilla files and test on-device.
- No resolution of external namespaces (`common.*`, `server_form.*`, Leaf/StarLib/Cherry are pack-specific, not engine APIs).
- `modifications`, variables, `$`-prefixed template args are preserved, not evaluated.
- Tabbed output and form labels/values require behavior-pack `server-ui` script logic.

## Layout

```
main.go
internal/app/app.go              # Run(), signal handling
internal/mcpserver/server.go     # tool + resource registration
internal/mcpserver/handlers.go   # tool handlers, install guidance
internal/bedrockui/validator.go, schema.go  # JSON-schema validation
internal/bedrockui/generators.go # action/modal/tabbed/grid + Cherry pack
internal/bedrockui/analyzer.go   # binding analysis + inheritance expansion
internal/bedrockui/reference.go + *.txt     # embedded docs
addon_scaffolding_instructions.txt          # new-addon + deploy.cjs conventions
```
