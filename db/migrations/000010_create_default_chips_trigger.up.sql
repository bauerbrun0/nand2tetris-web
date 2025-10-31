CREATE OR REPLACE FUNCTION insert_default_chips()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO chips (project_id, name, hdl)
  VALUES (
    NEW.id,
    'NotChip',
    '// A custom Not chip
CHIP NotChip {
    IN in;
    OUT out;

    PARTS:
    Nand(a = in, b = in, out = out);
}'
  );

  INSERT INTO chips (project_id, name, hdl)
  VALUES (
    NEW.id,
    'AndChip',
    '// A starter AndChip which uses a custom Not chip (NotChip)
// and a built-in Nand gate to implement the AND function.
CHIP AndChip {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a = a, b = b, out = nandOut);
    NotChip(in = nandOut, out = out);
}'
  );

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_projects_default_chips
AFTER INSERT ON projects
FOR EACH ROW
EXECUTE FUNCTION insert_default_chips();
