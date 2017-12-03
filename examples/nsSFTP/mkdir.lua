local sftp = require("nsSFTP")

function main()
    local destination = {
        HostPort = "localhost:22",
        User = "test",
        Password = "wireless"
    }

    local connection, err = sftp.connect(destination)
    if err ~= nil then
        error(err)
    end

    err = connection:mkdir("/tmp/hello")
    if err ~= nil then
        error(err)
    end

    err = connection:disconnect()
    if err ~= nil then
        error(err)
    end
end