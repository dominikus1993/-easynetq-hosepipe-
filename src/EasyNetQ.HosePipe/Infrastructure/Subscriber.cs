using RabbitMQ.Client;

namespace EasyNetQ.HosePipe.Infrastructure;

public sealed class Subscriber : BackgroundService
{
    private readonly IModel _subscriptionModel;
    private readonly IModel _publisherModel;
    public Subscriber(IConnection connection)
    {
        _subscriptionModel = connection.CreateModel();
        _publisherModel = connection.CreateModel();
    }

    protected override Task ExecuteAsync(CancellationToken stoppingToken)
    {
        throw new NotImplementedException();
    }

    public override void Dispose()
    {
        _subscriptionModel.Close();
        _subscriptionModel.Dispose();
        _publisherModel.Close();
        _publisherModel.Dispose();
        GC.SuppressFinalize(this);
    }
}