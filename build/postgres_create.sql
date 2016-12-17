CREATE OR REPLACE FUNCTION new_torrent_added() RETURNS TRIGGER AS $$
BEGIN
  PERFORM pg_notify(
      'new_torrent_added',

      json_build_object(
          'infoHash', NEW.info_hash,
          'name', NEW.name
      )::text
  );

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS new_torrent_added ON white_torrents;
CREATE TRIGGER new_torrent_added BEFORE INSERT ON white_torrents FOR EACH ROW EXECUTE PROCEDURE new_torrent_added();
