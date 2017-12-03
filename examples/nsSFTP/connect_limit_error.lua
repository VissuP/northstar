local sftp = require("nsSFTP")

function main()
    local destination = {
        HostPort = "localhost:22",
        User = "test",
        Password = "wireless"
    }

    for i = 1, 4, 1 do
        local _, err = sftp.connect(destination)
        if err ~= nil then
            error(err)
        end
    end
end