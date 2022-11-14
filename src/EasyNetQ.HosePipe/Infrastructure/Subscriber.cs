using RabbitMQ.Client;

namespace EasyNetQ.HosePipe.Infrastructure;

public class Subscriber : BackgroundService
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

    protected virtual void Dispose(bool disposing)
    {
        if (disposing)
        {
            _subscriptionModel.Close();
            _subscriptionModel.Dispose();
            _publisherModel.Close();
            _publisherModel.Dispose();
        }
    }

    public sealed override void Dispose()
    {
        Dispose(true);
        GC.SuppressFinalize(this);
    }
}