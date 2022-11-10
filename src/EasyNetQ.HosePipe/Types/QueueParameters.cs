namespace EasyNetQ.HosePipe.Types;

public sealed class QueueParameters
{
    public required string HostName { get; set; }
    public required int HostPort { get; set; }
    public required string VHost { get; set; }
    public required string Username { get; set; }
    public required string Password { get; set; }
    public required string QueueName { get; set; }
    public required bool Purge { get; set; }
    public required int NumberOfMessagesToRetrieve { get; set; }
    public required string MessagesOutputDirectory { get; set; }
    public TimeSpan ConfirmsTimeout { get; } = TimeSpan.FromSeconds(30);

    public QueueParameters()
    {
        // set some defaults
        HostName = "localhost";
        HostPort = -1;
        VHost = "/";
        Username = "guest";
        Password = "guest";
        Purge = false;
        NumberOfMessagesToRetrieve = 1000;
        MessagesOutputDirectory = Directory.GetCurrentDirectory();
    }
}